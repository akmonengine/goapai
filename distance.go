package goapai

// Distance calculates the distance between the current state value and a condition target value
// Returns 0 if the condition is already satisfied, otherwise returns the numeric distance
func (state State[T]) Distance(condition ConditionInterface) float32 {
	// Check if the condition key matches
	if state.Key != condition.GetKey() {
		return 0
	}

	// Handle different condition types
	switch cond := condition.(type) {
	case *Condition[int8]:
		if v, ok := any(state.Value).(int8); ok {
			return calculateNumericDistance(float64(v), float64(cond.Value), cond.Operator)
		}
	case *Condition[int]:
		if v, ok := any(state.Value).(int); ok {
			return calculateNumericDistance(float64(v), float64(cond.Value), cond.Operator)
		}
	case *Condition[uint8]:
		if v, ok := any(state.Value).(uint8); ok {
			return calculateNumericDistance(float64(v), float64(cond.Value), cond.Operator)
		}
	case *Condition[uint64]:
		if v, ok := any(state.Value).(uint64); ok {
			return calculateNumericDistance(float64(v), float64(cond.Value), cond.Operator)
		}
	case *Condition[float64]:
		if v, ok := any(state.Value).(float64); ok {
			return calculateNumericDistance(v, cond.Value, cond.Operator)
		}
	case *ConditionBool:
		if v, ok := any(state.Value).(bool); ok {
			if cond.Operator == EQUAL {
				if v == cond.Value {
					return 0
				}
				return 1
			} else if cond.Operator == NOT_EQUAL {
				if v != cond.Value {
					return 0
				}
				return 1
			}
		}
	case *ConditionString:
		if v, ok := any(state.Value).(string); ok {
			if cond.Operator == EQUAL {
				if v == cond.Value {
					return 0
				}
				return 1
			} else if cond.Operator == NOT_EQUAL {
				if v != cond.Value {
					return 0
				}
				return 1
			}
		}
	}

	return 0
}

// calculateNumericDistance computes the distance for numeric conditions based on operator
func calculateNumericDistance(current, target float64, op operator) float32 {
	switch op {
	case EQUAL:
		if current < target {
			return float32(target - current)
		}
		return float32(current - target)
	case NOT_EQUAL:
		if current == target {
			return 1.0
		}
		return 0.0
	case UPPER_OR_EQUAL:
		if current < target {
			return float32(target - current)
		}
		return 0.0
	case UPPER:
		if current <= target {
			return float32(target - current + 1)
		}
		return 0.0
	case LOWER_OR_EQUAL:
		if current > target {
			return float32(current - target)
		}
		return 0.0
	case LOWER:
		if current >= target {
			return float32(current - target + 1)
		}
		return 0.0
	}
	return 0.0
}
