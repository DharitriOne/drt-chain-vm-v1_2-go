package scenarioexec

// ProtectedKeyPrefix prefixes all reserved storage. Only the protocol can write to keys starting with this.
const ProtectedKeyPrefix = "N" + "U" + "M" + "B" + "A" + "T"

// RewardKey is the storage key where the protocol writes when sending out rewards.
const RewardKey = ProtectedKeyPrefix + "reward"
