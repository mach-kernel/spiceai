package pods

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"testing"
	"time"

	"github.com/bradleyjkemp/cupaloy"
	"github.com/spiceai/spice/pkg/csv"
	"github.com/stretchr/testify/assert"
)

var snapshotter = cupaloy.New(cupaloy.SnapshotSubdirectory("../../test/assets/snapshots/pods"))

func TestPod(t *testing.T) {
	manifestsToTest := []string{"trader.yaml", "trader-infer.yaml", "cartpole-v1.yaml"}

	for _, manifestToTest := range manifestsToTest {
		manifestPath := filepath.Join("../../test/assets/pods/manifests", manifestToTest)

		pod, err := LoadPodFromManifest(manifestPath)
		if err != nil {
			t.Error(err)
			return
		}

		t.Run(fmt.Sprintf("Base Properties - %s", manifestToTest), testBasePropertiesFunc(pod))
		t.Run(fmt.Sprintf("FieldNames() - %s", manifestToTest), testFieldNamesFunc(pod))
		t.Run(fmt.Sprintf("Rewards() - %s", manifestToTest), testRewardsFunc(pod))
		t.Run(fmt.Sprintf("Actions() - %s", manifestToTest), testActionsFunc(pod))
		t.Run(fmt.Sprintf("CachedCsv() - %s", manifestToTest), testCachedCsvFunc(pod))
		t.Run(fmt.Sprintf("AddLocalState() - %s", manifestToTest), testAddLocalStateFunc(pod))
		t.Run(fmt.Sprintf("AddLocalState()/CachedCsv() - %s", manifestToTest), testAddLocalStateCachedCsvFunc(pod))
	}
}

// Tests base properties
func testBasePropertiesFunc(pod *Pod) func(*testing.T) {
	return func(t *testing.T) {

		actual := pod.Hash()

		var expected string

		switch pod.Name {
		case "trader":
			expected = "d2cc6c526c3f48b630a251fb1b586cf4"
		case "trader-infer":
			expected = "95c0e9b0b65ae7b50af26b4582126624"
		case "cartpole-v1":
			expected = "f739fd1bf82dc26a864dca4778a43331"
		}

		assert.Equal(t, expected, actual, "invalid pod.Hash()")

		actual = pod.ManifestPath()

		switch pod.Name {
		case "trader":
			expected = "../../test/assets/pods/manifests/trader.yaml"
		case "trader-infer":
			expected = "../../test/assets/pods/manifests/trader-infer.yaml"
		case "cartpole-v1":
			expected = "../../test/assets/pods/manifests/cartpole-v1.yaml"
		}

		assert.Equal(t, expected, actual, "invalid pod.ManifestPath()")

		actual = fmt.Sprintf("%d", pod.Epoch().Unix())

		switch pod.Name {
		case "trader":
			expected = "1605312000"
		case "trader-infer":
			actual = actual[:8] // Reduce precision to test
			expected = fmt.Sprintf("%d", time.Now().Add(-pod.Period()).Unix())[:8]
		case "cartpole-v1":
			actual = actual[:8] // Reduce precision to test
			expected = fmt.Sprintf("%d", time.Now().Add(-pod.Period()).Unix())[:8]
		}

		assert.Equal(t, expected, actual, "invalid pod.Epoch()")

		actual = pod.Period().String()

		switch pod.Name {
		case "trader":
			expected = "17h0m0s"
		case "trader-infer":
			expected = "72h0m0s"
		case "cartpole-v1":
			expected = "50m0s"
		}

		assert.Equal(t, expected, actual, "invalid pod.Period()")

		actual = pod.Interval().String()

		switch pod.Name {
		case "trader":
			expected = "17m0s"
		case "trader-infer":
			expected = "1m0s"
		case "cartpole-v1":
			expected = "10s"
		}

		assert.Equal(t, expected, actual, "invalid pod.Interval()")

		actual = pod.Granularity().String()

		switch pod.Name {
		case "trader":
			expected = "17s"
		case "trader-infer":
			expected = "10s"
		case "cartpole-v1":
			expected = "10s"
		}

		assert.Equal(t, expected, actual, "invalid pod.Granularity()")
	}
}

// Tests FieldNames() getter
func testFieldNamesFunc(pod *Pod) func(*testing.T) {
	return func(t *testing.T) {
		actual := pod.FieldNames()

		var expected []string

		switch pod.Name {
		case "trader":
			fallthrough
		case "trader-infer":
			expected = []string{
				"coinbase.btcusd.close",
				"local.portfolio.btc_balance",
				"local.portfolio.usd_balance",
			}
		case "cartpole-v1":
			expected = []string{
				"gym.CartPole.cart_position",
				"gym.CartPole.cart_velocity",
				"gym.CartPole.is_done",
				"gym.CartPole.pole_angle",
				"gym.CartPole.pole_angular_velocity",
			}
		}

		assert.Equal(t, expected, actual, "invalid pod.FieldNames()")
	}
}

// Tests Rewards() getter
func testRewardsFunc(pod *Pod) func(*testing.T) {
	return func(t *testing.T) {
		actual := pod.Rewards()

		var expected map[string]string

		switch pod.Name {
		case "trader":
			expected = map[string]string{
				"buy":  "reward = 1",
				"sell": "reward = 1",
				"hold": "reward = 1",
			}
		case "trader-infer":
			expected = map[string]string{
				"buy":  "reward = 1",
				"sell": "reward = 1",
			}
		case "cartpole-v1":
			expected = map[string]string{
				"left":  "reward = 0 if new_state.gym_CartPole_is_done == 1 else 1",
				"right": "reward = 0 if new_state.gym_CartPole_is_done == 1 else 1",
			}
		}

		assert.Equal(t, expected, actual, "invalid pod.Rewards()")
	}
}

// Tests Actions() getter
func testActionsFunc(pod *Pod) func(*testing.T) {
	return func(t *testing.T) {
		actual := pod.Actions()

		var expected map[string]string

		switch pod.Name {
		case "trader":
			expected = map[string]string{
				"buy":  "local.portfolio.usd_balance -= coinbase.btcusd.close\nlocal.portfolio.btc_balance += 1",
				"hold": "",
				"sell": "local.portfolio.usd_balance += coinbase.btcusd.close\nlocal.portfolio.btc_balance -= 1",
			}
		case "trader-infer":
			expected = map[string]string{
				"buy":  "local.portfolio.usd_balance -= args.price\nlocal.portfolio.btc_balance += 1",
				"sell": "local.portfolio.usd_balance += args.price\nlocal.portfolio.btc_balance -= 1",
			}
		case "cartpole-v1":
			expected = map[string]string{
				"left":  "passthru",
				"right": "passthru",
			}
		default:
			t.Errorf("invalid pod %s", pod.Name)
		}

		assert.Equal(t, expected, actual, "invalid pod.Actions()")
	}
}

// Tests CachedCsv() getter
func testCachedCsvFunc(pod *Pod) func(*testing.T) {
	return func(t *testing.T) {
		_, err := pod.FetchNewData()
		if err != nil {
			t.Error(err)
			return
		}

		actual := pod.CachedCsv()

		snapshotter.SnapshotT(t, actual)
	}
}

// Tests AddLocalState()
func testAddLocalStateFunc(pod *Pod) func(*testing.T) {
	return func(t *testing.T) {
		data, err := ioutil.ReadFile("../../test/assets/data/csv/trader_input.csv")
		if err != nil {
			t.Error(err)
			return
		}

		reader := bytes.NewReader(data)
		newState, err := csv.ProcessCsvByPath(reader, nil)
		if err != nil {
			t.Error(err)
		}

		pod.AddLocalState(newState...)
	}
}

// Tests AddLocalStateCachedCsv()
func testAddLocalStateCachedCsvFunc(pod *Pod) func(*testing.T) {
	return func(t *testing.T) {
		data, err := ioutil.ReadFile("../../test/assets/data/csv/trader_input.csv")
		if err != nil {
			t.Error(err)
			return
		}

		reader := bytes.NewReader(data)
		newState, err := csv.ProcessCsvByPath(reader, nil)
		if err != nil {
			t.Error(err)
		}

		assert.Equal(t, 2, len(newState), "expected two state objects, one for local and one for coinbase")

		pod.AddLocalState(newState...)

		actual := pod.CachedCsv()

		snapshotter.SnapshotT(t, actual)
	}
}