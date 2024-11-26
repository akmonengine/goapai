package benchmark

import (
	"fmt"
	"goapai"
	"os"
	"runtime/pprof"
	"runtime/trace"
	"testing"
	"time"
)

const (
	ATTRIBUTE_HUNGRY = iota
	ATTRIBUTE_1
	ATTRIBUTE_2
	ATTRIBUTE_3
	ATTRIBUTE_4
	ATTRIBUTE_5
	ATTRIBUTE_6
	ATTRIBUTE_7
	ATTRIBUTE_8
	ATTRIBUTE_9
	ATTRIBUTE_10
	ATTRIBUTE_11
	ATTRIBUTE_12
	ATTRIBUTE_13
	ATTRIBUTE_14
	ATTRIBUTE_15
	ATTRIBUTE_16
	ATTRIBUTE_17
	ATTRIBUTE_18
	ATTRIBUTE_19
	ATTRIBUTE_20
	ATTRIBUTE_21
	ATTRIBUTE_22
	ATTRIBUTE_23
	ATTRIBUTE_24
	ATTRIBUTE_25
	ATTRIBUTE_26
	ATTRIBUTE_27
	ATTRIBUTE_28
	ATTRIBUTE_29
	ATTRIBUTE_30
	ATTRIBUTE_31
	ATTRIBUTE_32
	ATTRIBUTE_33
	ATTRIBUTE_34
	ATTRIBUTE_35
	ATTRIBUTE_36
	ATTRIBUTE_37
	ATTRIBUTE_38
	ATTRIBUTE_39
	ATTRIBUTE_40
	ATTRIBUTE_41
	ATTRIBUTE_42
	ATTRIBUTE_43
	ATTRIBUTE_44
	ATTRIBUTE_45
	ATTRIBUTE_46
	ATTRIBUTE_47
	ATTRIBUTE_48
	ATTRIBUTE_49
	ATTRIBUTE_50
	ATTRIBUTE_51
	ATTRIBUTE_52
	ATTRIBUTE_53
	ATTRIBUTE_54
	ATTRIBUTE_55
	ATTRIBUTE_56
	ATTRIBUTE_57
	ATTRIBUTE_58
	ATTRIBUTE_59
	ATTRIBUTE_60
	ATTRIBUTE_61
	ATTRIBUTE_62
	ATTRIBUTE_63
	ATTRIBUTE_64
	ATTRIBUTE_65
	ATTRIBUTE_66
	ATTRIBUTE_67
	ATTRIBUTE_68
	ATTRIBUTE_69
	ATTRIBUTE_70
	ATTRIBUTE_71
	ATTRIBUTE_72
	ATTRIBUTE_73
	ATTRIBUTE_74
	ATTRIBUTE_75
	ATTRIBUTE_76
	ATTRIBUTE_77
	ATTRIBUTE_78
	ATTRIBUTE_79
	ATTRIBUTE_80
	ATTRIBUTE_81
	ATTRIBUTE_82
	ATTRIBUTE_83
	ATTRIBUTE_84
	ATTRIBUTE_85
	ATTRIBUTE_86
	ATTRIBUTE_87
	ATTRIBUTE_88
	ATTRIBUTE_89
	ATTRIBUTE_90
	ATTRIBUTE_91
	ATTRIBUTE_92
	ATTRIBUTE_93
	ATTRIBUTE_94
	ATTRIBUTE_95
	ATTRIBUTE_96
	ATTRIBUTE_97
	ATTRIBUTE_98
	ATTRIBUTE_99
	ATTRIBUTE_100

	SEE_APPLE
	HAVE_APPLE
	HAVE_MEAL
	EAT
)

type Attributes struct {
	happiness  int
	hungry     bool
	hungriness int
	sleepy     bool
	attribute  bool
}
type Entity struct {
	agent      goapai.Agent
	attributes Attributes
}

func BenchmarkGoapAI(b *testing.B) {
	b.StopTimer()

	actions := goapai.Actions{}
	//goapai.AddAction("action 1", 1, false, goapai.Conditions{
	//	goapai.ConditionBool{Key: ATTRIBUTE_1, Value: false},
	//}, goapai.Effect{data: map[goapai.StateKey]goapai.StateInterface{ATTRIBUTE_2: goapai.StateBool{Value: false}}}),
	//
	//goapai.AddAction("action 2", 2, false, goapai.Conditions{
	//	goapai.ConditionBool{Key: ATTRIBUTE_1, Value: false},
	//}, goapai.Effect{data: map[goapai.StateKey]goapai.StateInterface{ATTRIBUTE_3: goapai.StateBool{Value: false}}}),
	//
	//goapai.AddAction("action 3", 3, false, goapai.Conditions{
	//	goapai.ConditionBool{Key: ATTRIBUTE_1, Value: false},
	//}, goapai.Effect{data: map[goapai.StateKey]goapai.StateInterface{ATTRIBUTE_4: goapai.StateBool{Value: false}}}),
	//
	//goapai.AddAction("action 4", 4, false, goapai.Conditions{
	//	goapai.ConditionBool{Key: ATTRIBUTE_1, Value: false},
	//}, goapai.Effect{data: map[goapai.StateKey]goapai.StateInterface{ATTRIBUTE_5: goapai.StateBool{Value: false}}}),
	//
	//goapai.AddAction("action 5", 5, false, goapai.Conditions{
	//	goapai.ConditionBool{Key: ATTRIBUTE_1, Value: false},
	//}, goapai.Effect{data: map[goapai.StateKey]goapai.StateInterface{ATTRIBUTE_6: goapai.StateBool{Value: false}}}),
	//
	//goapai.AddAction("action 6", 6, false, goapai.Conditions{
	//	goapai.ConditionBool{Key: ATTRIBUTE_1, Value: false},
	//}, goapai.Effect{data: map[goapai.StateKey]goapai.StateInterface{ATTRIBUTE_7: goapai.StateBool{Value: false}}}),
	//
	//goapai.AddAction("action 7", 7, false, goapai.Conditions{
	//	goapai.ConditionBool{Key: ATTRIBUTE_1, Value: false},
	//}, goapai.Effect{data: map[goapai.StateKey]goapai.StateInterface{ATTRIBUTE_8: goapai.StateBool{Value: false}}}),
	//
	//goapai.AddAction("action 8", 8, false, goapai.Conditions{
	//	goapai.ConditionBool{Key: ATTRIBUTE_1, Value: false},
	//}, goapai.Effect{data: map[goapai.StateKey]goapai.StateInterface{ATTRIBUTE_9: goapai.StateBool{Value: false}}}),
	//
	//goapai.AddAction("action 9", 8, false, goapai.Conditions{
	//	goapai.ConditionBool{Key: ATTRIBUTE_1, Value: false},
	//}, goapai.Effect{data: map[goapai.StateKey]goapai.StateInterface{ATTRIBUTE_9: goapai.StateBool{Value: false}}}),
	//
	//goapai.AddAction("action 10", 10, false, goapai.Conditions{
	//	goapai.ConditionBool{Key: ATTRIBUTE_1, Value: true},
	//}, goapai.Effect{data: map[goapai.StateKey]goapai.StateInterface{ATTRIBUTE_11: goapai.StateBool{Value: false}}}),
	//
	//goapai.AddAction("action 11", 11, false, goapai.Conditions{
	//	goapai.ConditionBool{Key: ATTRIBUTE_1, Value: true},
	//}, goapai.Effect{data: map[goapai.StateKey]goapai.StateInterface{ATTRIBUTE_12: goapai.StateBool{Value: false}}}),
	//
	//goapai.AddAction("action 12", 12, false, goapai.Conditions{
	//	goapai.ConditionBool{Key: ATTRIBUTE_1, Value: true},
	//}, goapai.Effect{data: map[goapai.StateKey]goapai.StateInterface{ATTRIBUTE_13: goapai.StateBool{Value: false}}}),
	//
	//goapai.AddAction("action 13", 13, false, goapai.Conditions{
	//	goapai.ConditionBool{Key: ATTRIBUTE_1, Value: true},
	//}, goapai.Effect{data: map[goapai.StateKey]goapai.StateInterface{ATTRIBUTE_14: goapai.StateBool{Value: false}}}),
	//
	//goapai.AddAction("action 14", 14, false, goapai.Conditions{
	//	goapai.ConditionBool{Key: ATTRIBUTE_1, Value: true},
	//}, goapai.Effect{data: map[goapai.StateKey]goapai.StateInterface{ATTRIBUTE_15: goapai.StateBool{Value: false}}}),
	//
	//goapai.AddAction("action 15", 15, false, goapai.Conditions{
	//	goapai.ConditionBool{Key: ATTRIBUTE_1, Value: true},
	//}, goapai.Effect{data: map[goapai.StateKey]goapai.StateInterface{ATTRIBUTE_16: goapai.StateBool{Value: false}}}),
	//
	//goapai.AddAction("action 16", 16, false, goapai.Conditions{
	//	goapai.ConditionBool{Key: ATTRIBUTE_1, Value: true},
	//}, goapai.Effect{data: map[goapai.StateKey]goapai.StateInterface{ATTRIBUTE_17: goapai.StateBool{Value: false}}}),
	//
	//goapai.AddAction("action 17", 17, false, goapai.Conditions{
	//	goapai.ConditionBool{Key: ATTRIBUTE_1, Value: true},
	//}, goapai.Effect{data: map[goapai.StateKey]goapai.StateInterface{ATTRIBUTE_18: goapai.StateBool{Value: false}}}),
	//
	//goapai.AddAction("action 18", 18, false, goapai.Conditions{
	//	goapai.ConditionBool{Key: ATTRIBUTE_1, Value: true},
	//}, goapai.Effect{data: map[goapai.StateKey]goapai.StateInterface{ATTRIBUTE_19: goapai.StateBool{Value: false}}}),
	//
	//goapai.AddAction("action 19", 19, false, goapai.Conditions{
	//	goapai.ConditionBool{Key: ATTRIBUTE_1, Value: true},
	//}, goapai.Effect{data: map[goapai.StateKey]goapai.StateInterface{ATTRIBUTE_20: goapai.StateBool{Value: false}}}),
	//
	//goapai.AddAction("action 20", 20, false, goapai.Conditions{
	//	goapai.ConditionBool{Key: ATTRIBUTE_1, Value: true},
	//}, goapai.Effect{data: map[goapai.StateKey]goapai.StateInterface{}}),
	//
	//goapai.AddAction("action 21", 21, false, goapai.Conditions{
	//	goapai.ConditionBool{Key: ATTRIBUTE_1, Value: true},
	//}, goapai.Effect{data: map[goapai.StateKey]goapai.StateInterface{ATTRIBUTE_22: goapai.StateBool{Value: false}}}),
	//
	//goapai.AddAction("action 22", 22, false, goapai.Conditions{
	//	goapai.ConditionBool{Key: ATTRIBUTE_1, Value: true},
	//}, goapai.Effect{data: map[goapai.StateKey]goapai.StateInterface{ATTRIBUTE_23: goapai.StateBool{Value: false}}}),
	//
	//goapai.AddAction("action 23", 23, false, goapai.Conditions{
	//	goapai.ConditionBool{Key: ATTRIBUTE_1, Value: true},
	//}, goapai.Effect{data: map[goapai.StateKey]goapai.StateInterface{ATTRIBUTE_24: goapai.StateBool{Value: false}}}),
	//
	//goapai.AddAction("action 24", 24, false, goapai.Conditions{
	//	goapai.ConditionBool{Key: ATTRIBUTE_1, Value: true},
	//}, goapai.Effect{data: map[goapai.StateKey]goapai.StateInterface{ATTRIBUTE_25: goapai.StateBool{Value: false}}}),
	//
	//goapai.AddAction("action 25", 25, false, goapai.Conditions{
	//	goapai.ConditionBool{Key: ATTRIBUTE_1, Value: true},
	//}, goapai.Effect{data: map[goapai.StateKey]goapai.StateInterface{ATTRIBUTE_26: goapai.StateBool{Value: false}}}),
	//
	//goapai.AddAction("action 26", 26, false, goapai.Conditions{
	//	goapai.ConditionBool{Key: ATTRIBUTE_1, Value: true},
	//}, goapai.Effect{data: map[goapai.StateKey]goapai.StateInterface{ATTRIBUTE_27: goapai.StateBool{Value: false}}}),
	//
	//goapai.AddAction("action 27", 27, false, goapai.Conditions{
	//	goapai.ConditionBool{Key: ATTRIBUTE_1, Value: true},
	//}, goapai.Effect{data: map[goapai.StateKey]goapai.StateInterface{ATTRIBUTE_28: goapai.StateBool{Value: false}}}),
	//
	//goapai.AddAction("action 28", 28, false, goapai.Conditions{
	//	goapai.ConditionBool{Key: ATTRIBUTE_1, Value: true},
	//}, goapai.Effect{data: map[goapai.StateKey]goapai.StateInterface{ATTRIBUTE_29: goapai.StateBool{Value: false}}}),
	//
	//goapai.AddAction("action 29", 29, false, goapai.Conditions{
	//	goapai.ConditionBool{Key: ATTRIBUTE_1, Value: true},
	//}, goapai.Effect{data: map[goapai.StateKey]goapai.StateInterface{ATTRIBUTE_30: goapai.StateBool{Value: false}}}),
	//
	//goapai.AddAction("action 30", 30, false, goapai.Conditions{
	//	goapai.ConditionBool{Key: ATTRIBUTE_1, Value: true},
	//}, goapai.Effect{data: map[goapai.StateKey]goapai.StateInterface{}}),
	//
	//goapai.AddAction("action 31", 31, false, goapai.Conditions{
	//	goapai.ConditionBool{Key: ATTRIBUTE_1, Value: true},
	//}, goapai.Effect{data: map[goapai.StateKey]goapai.StateInterface{ATTRIBUTE_32: goapai.StateBool{Value: false}}}),
	//
	//goapai.AddAction("action 32", 32, false, goapai.Conditions{
	//	goapai.ConditionBool{Key: ATTRIBUTE_1, Value: true},
	//}, goapai.Effect{data: map[goapai.StateKey]goapai.StateInterface{ATTRIBUTE_33: goapai.StateBool{Value: false}}}),
	//
	//goapai.AddAction("action 33", 33, false, goapai.Conditions{
	//	goapai.ConditionBool{Key: ATTRIBUTE_1, Value: true},
	//}, goapai.Effect{data: map[goapai.StateKey]goapai.StateInterface{}}),
	//
	//goapai.AddAction("action 34", 34, false, goapai.Conditions{
	//	goapai.ConditionBool{Key: ATTRIBUTE_1, Value: true},
	//}, goapai.Effect{data: map[goapai.StateKey]goapai.StateInterface{ATTRIBUTE_35: goapai.StateBool{Value: false}}}),
	//
	//goapai.AddAction("action 35", 35, false, goapai.Conditions{
	//	goapai.ConditionBool{Key: ATTRIBUTE_1, Value: true},
	//}, goapai.Effect{data: map[goapai.StateKey]goapai.StateInterface{ATTRIBUTE_36: goapai.StateBool{Value: false}}}),
	//
	//goapai.AddAction("action 36", 36, false, goapai.Conditions{
	//	goapai.ConditionBool{Key: ATTRIBUTE_1, Value: true},
	//}, goapai.Effect{data: map[goapai.StateKey]goapai.StateInterface{}}),
	//
	//goapai.AddAction("action 37", 37, false, goapai.Conditions{
	//	goapai.ConditionBool{Key: ATTRIBUTE_1, Value: true},
	//}, goapai.Effect{data: map[goapai.StateKey]goapai.StateInterface{ATTRIBUTE_38: goapai.StateBool{Value: false}}}),
	//
	//goapai.AddAction("action 38", 38, false, goapai.Conditions{
	//	goapai.ConditionBool{Key: ATTRIBUTE_1, Value: true},
	//}, goapai.Effect{data: map[goapai.StateKey]goapai.StateInterface{ATTRIBUTE_39: goapai.StateBool{Value: false}}}),
	//
	//goapai.AddAction("action 39", 39, false, goapai.Conditions{
	//	goapai.ConditionBool{Key: ATTRIBUTE_1, Value: true},
	//}, goapai.Effect{data: map[goapai.StateKey]goapai.StateInterface{}}),
	//
	//goapai.AddAction("action 40", 40, false, goapai.Conditions{
	//	goapai.ConditionBool{Key: ATTRIBUTE_1, Value: true},
	//}, goapai.Effect{data: map[goapai.StateKey]goapai.StateInterface{ATTRIBUTE_41: goapai.StateBool{Value: false}}}),
	//
	//goapai.AddAction("action 41", 41, false, goapai.Conditions{
	//	goapai.ConditionBool{Key: ATTRIBUTE_1, Value: true},
	//}, goapai.Effect{data: map[goapai.StateKey]goapai.StateInterface{ATTRIBUTE_42: goapai.StateBool{Value: false}}}),
	//
	//goapai.AddAction("action 42", 42, false, goapai.Conditions{
	//	goapai.ConditionBool{Key: ATTRIBUTE_1, Value: true},
	//}, goapai.Effect{data: map[goapai.StateKey]goapai.StateInterface{}}),
	//
	//goapai.AddAction("action 43", 43, false, goapai.Conditions{
	//	goapai.ConditionBool{Key: ATTRIBUTE_1, Value: true},
	//}, goapai.Effect{data: map[goapai.StateKey]goapai.StateInterface{ATTRIBUTE_44: goapai.StateBool{Value: false}}}),
	//
	//goapai.AddAction("action 44", 44, false, goapai.Conditions{
	//	goapai.ConditionBool{Key: ATTRIBUTE_1, Value: true},
	//}, goapai.Effect{data: map[goapai.StateKey]goapai.StateInterface{ATTRIBUTE_45: goapai.StateBool{Value: false}}}),
	//
	//goapai.AddAction("action 45", 45, false, goapai.Conditions{
	//	goapai.ConditionBool{Key: ATTRIBUTE_1, Value: true},
	//}, goapai.Effect{data: map[goapai.StateKey]goapai.StateInterface{}}),
	//
	//goapai.AddAction("action 46", 46, false, goapai.Conditions{
	//	goapai.ConditionBool{Key: ATTRIBUTE_1, Value: true},
	//}, goapai.Effect{data: map[goapai.StateKey]goapai.StateInterface{ATTRIBUTE_47: goapai.StateBool{Value: false}}}),
	//
	//goapai.AddAction("action 47", 47, false, goapai.Conditions{
	//	goapai.ConditionBool{Key: ATTRIBUTE_1, Value: true},
	//}, goapai.Effect{data: map[goapai.StateKey]goapai.StateInterface{ATTRIBUTE_48: goapai.StateBool{Value: false}}}),
	//
	//goapai.AddAction("action 48", 48, false, goapai.Conditions{
	//	goapai.ConditionBool{Key: ATTRIBUTE_1, Value: true},
	//}, goapai.Effect{data: map[goapai.StateKey]goapai.StateInterface{}}),
	//
	//goapai.AddAction("action 49", 49, false, goapai.Conditions{
	//	goapai.ConditionBool{Key: ATTRIBUTE_1, Value: true},
	//}, goapai.Effect{data: map[goapai.StateKey]goapai.StateInterface{ATTRIBUTE_50: goapai.StateBool{Value: false}}}),
	//
	//goapai.AddAction("action 50", 50, false, goapai.Conditions{
	//	goapai.ConditionBool{Key: ATTRIBUTE_1, Value: true},
	//}, goapai.Effect{data: map[goapai.StateKey]goapai.StateInterface{}}),
	//
	//goapai.AddAction("action 51", 51, false, goapai.Conditions{
	//	goapai.ConditionBool{Key: ATTRIBUTE_1, Value: true},
	//}, goapai.Effect{data: map[goapai.StateKey]goapai.StateInterface{}}),
	//
	//goapai.AddAction("action 52", 52, false, goapai.Conditions{
	//	goapai.ConditionBool{Key: ATTRIBUTE_1, Value: true},
	//}, goapai.Effect{data: map[goapai.StateKey]goapai.StateInterface{ATTRIBUTE_53: goapai.StateBool{Value: false}}}),
	//
	//goapai.AddAction("action 53", 53, false, goapai.Conditions{
	//	goapai.ConditionBool{Key: ATTRIBUTE_1, Value: true},
	//}, goapai.Effect{data: map[goapai.StateKey]goapai.StateInterface{ATTRIBUTE_54: goapai.StateBool{Value: false}}}),
	//
	//goapai.AddAction("action 54", 54, false, goapai.Conditions{
	//	goapai.ConditionBool{Key: ATTRIBUTE_1, Value: true},
	//}, goapai.Effect{data: map[goapai.StateKey]goapai.StateInterface{}}),
	//
	//goapai.AddAction("action 55", 55, false, goapai.Conditions{
	//	goapai.ConditionBool{Key: ATTRIBUTE_1, Value: true},
	//}, goapai.Effect{data: map[goapai.StateKey]goapai.StateInterface{ATTRIBUTE_56: goapai.StateBool{Value: false}}}),
	//
	//goapai.AddAction("action 56", 56, false, goapai.Conditions{
	//	goapai.ConditionBool{Key: ATTRIBUTE_1, Value: true},
	//}, goapai.Effect{data: map[goapai.StateKey]goapai.StateInterface{ATTRIBUTE_57: goapai.StateBool{Value: false}}}),
	//
	//goapai.AddAction("action 57", 57, false, goapai.Conditions{
	//	goapai.ConditionBool{Key: ATTRIBUTE_1, Value: true},
	//}, goapai.Effect{data: map[goapai.StateKey]goapai.StateInterface{}}),
	//
	//goapai.AddAction("action 58", 58, false, goapai.Conditions{
	//	goapai.ConditionBool{Key: ATTRIBUTE_1, Value: true},
	//}, goapai.Effect{data: map[goapai.StateKey]goapai.StateInterface{ATTRIBUTE_59: goapai.StateBool{Value: false}}}),
	//
	//goapai.AddAction("action 59", 59, false, goapai.Conditions{
	//	goapai.ConditionBool{Key: ATTRIBUTE_1, Value: true},
	//}, goapai.Effect{data: map[goapai.StateKey]goapai.StateInterface{ATTRIBUTE_60: goapai.StateBool{Value: false}}}),
	//
	//goapai.AddAction("action 60", 60, false, goapai.Conditions{
	//	goapai.ConditionBool{Key: ATTRIBUTE_1, Value: true},
	//}, goapai.Effect{data: map[goapai.StateKey]goapai.StateInterface{}}),
	//
	//goapai.AddAction("action 61", 61, false, goapai.Conditions{
	//	goapai.ConditionBool{Key: ATTRIBUTE_1, Value: true},
	//}, goapai.Effect{data: map[goapai.StateKey]goapai.StateInterface{ATTRIBUTE_62: goapai.StateBool{Value: false}}}),
	//
	//goapai.AddAction("action 62", 62, false, goapai.Conditions{
	//	goapai.ConditionBool{Key: ATTRIBUTE_1, Value: true},
	//}, goapai.Effect{data: map[goapai.StateKey]goapai.StateInterface{ATTRIBUTE_63: goapai.StateBool{Value: false}}}),
	//
	//goapai.AddAction("action 63", 63, false, goapai.Conditions{
	//	goapai.ConditionBool{Key: ATTRIBUTE_1, Value: true},
	//}, goapai.Effect{data: map[goapai.StateKey]goapai.StateInterface{}}),
	//
	//goapai.AddAction("action 64", 64, false, goapai.Conditions{
	//	goapai.ConditionBool{Key: ATTRIBUTE_1, Value: true},
	//}, goapai.Effect{data: map[goapai.StateKey]goapai.StateInterface{ATTRIBUTE_65: goapai.StateBool{Value: false}}}),
	//
	//goapai.AddAction("action 65", 65, false, goapai.Conditions{
	//	goapai.ConditionBool{Key: ATTRIBUTE_1, Value: true},
	//}, goapai.Effect{data: map[goapai.StateKey]goapai.StateInterface{ATTRIBUTE_66: goapai.StateBool{Value: false}}}),
	//
	//goapai.AddAction("action 66", 66, false, goapai.Conditions{
	//	goapai.ConditionBool{Key: ATTRIBUTE_1, Value: true},
	//}, goapai.Effect{data: map[goapai.StateKey]goapai.StateInterface{}}),
	//
	//goapai.AddAction("action 67", 67, false, goapai.Conditions{
	//	goapai.ConditionBool{Key: ATTRIBUTE_1, Value: true},
	//}, goapai.Effect{data: map[goapai.StateKey]goapai.StateInterface{ATTRIBUTE_68: goapai.StateBool{Value: false}}}),
	//
	//goapai.AddAction("action 68", 68, false, goapai.Conditions{
	//	goapai.ConditionBool{Key: ATTRIBUTE_1, Value: true},
	//}, goapai.Effect{data: map[goapai.StateKey]goapai.StateInterface{ATTRIBUTE_69: goapai.StateBool{Value: false}}}),
	//
	//goapai.AddAction("action 69", 69, false, goapai.Conditions{
	//	goapai.ConditionBool{Key: ATTRIBUTE_1, Value: true},
	//}, goapai.Effect{data: map[goapai.StateKey]goapai.StateInterface{}}),
	//
	//goapai.AddAction("action 70", 70, false, goapai.Conditions{
	//	goapai.ConditionBool{Key: ATTRIBUTE_1, Value: true},
	//}, goapai.Effect{data: map[goapai.StateKey]goapai.StateInterface{ATTRIBUTE_70: goapai.StateBool{Value: false}}}),
	//
	//goapai.AddAction("action 71", 71, false, goapai.Conditions{
	//	goapai.ConditionBool{Key: ATTRIBUTE_1, Value: true},
	//}, goapai.Effect{data: map[goapai.StateKey]goapai.StateInterface{}}),
	//
	//goapai.AddAction("action 72", 72, false, goapai.Conditions{
	//	goapai.ConditionBool{Key: ATTRIBUTE_1, Value: true},
	//}, goapai.Effect{data: map[goapai.StateKey]goapai.StateInterface{ATTRIBUTE_73: goapai.StateBool{Value: false}}}),
	//
	//goapai.AddAction("action 73", 73, false, goapai.Conditions{
	//	goapai.ConditionBool{Key: ATTRIBUTE_1, Value: true},
	//}, goapai.Effect{data: map[goapai.StateKey]goapai.StateInterface{ATTRIBUTE_74: goapai.StateBool{Value: false}}}),
	//
	//goapai.AddAction("action 74", 74, false, goapai.Conditions{
	//	goapai.ConditionBool{Key: ATTRIBUTE_1, Value: true},
	//}, goapai.Effect{data: map[goapai.StateKey]goapai.StateInterface{}}),
	//
	//goapai.AddAction("action 75", 75, false, goapai.Conditions{
	//	goapai.ConditionBool{Key: ATTRIBUTE_1, Value: true},
	//}, goapai.Effect{data: map[goapai.StateKey]goapai.StateInterface{ATTRIBUTE_76: goapai.StateBool{Value: false}}}),
	//
	//goapai.AddAction("action 76", 76, false, goapai.Conditions{
	//	goapai.ConditionBool{Key: ATTRIBUTE_1, Value: true},
	//}, goapai.Effect{data: map[goapai.StateKey]goapai.StateInterface{ATTRIBUTE_77: goapai.StateBool{Value: false}}}),
	//
	//goapai.AddAction("action 77", 77, false, goapai.Conditions{
	//	goapai.ConditionBool{Key: ATTRIBUTE_1, Value: true},
	//}, goapai.Effect{data: map[goapai.StateKey]goapai.StateInterface{}}),
	//
	//goapai.AddAction("action 78", 78, false, goapai.Conditions{
	//	goapai.ConditionBool{Key: ATTRIBUTE_1, Value: true},
	//}, goapai.Effect{data: map[goapai.StateKey]goapai.StateInterface{ATTRIBUTE_79: goapai.StateBool{Value: false}}}),
	//
	//goapai.AddAction("action 79", 79, false, goapai.Conditions{
	//	goapai.ConditionBool{Key: ATTRIBUTE_1, Value: true},
	//}, goapai.Effect{data: map[goapai.StateKey]goapai.StateInterface{ATTRIBUTE_80: goapai.StateBool{Value: false}}}),
	//
	//goapai.AddAction("action 80", 80, false, goapai.Conditions{
	//	goapai.ConditionBool{Key: ATTRIBUTE_1, Value: true},
	//}, goapai.Effect{data: map[goapai.StateKey]goapai.StateInterface{}}),

	actions.AddAction("action 81", 1, true, goapai.Conditions{
		&goapai.ConditionBool{Key: ATTRIBUTE_HUNGRY, Value: true},
	}, goapai.Effects{
		goapai.Effect[int]{
			Key:      ATTRIBUTE_3,
			Value:    1,
			Operator: goapai.EFFECT_ARITHMETIC_SET,
		},
		goapai.EffectString{
			Key:      ATTRIBUTE_4,
			Value:    "attribute4",
			Operator: goapai.EFFECT_ARITHMETIC_ADD,
		},
		goapai.Effect[int]{
			Key:      ATTRIBUTE_6,
			Value:    1,
			Operator: goapai.EFFECT_ARITHMETIC_ADD,
		},
	})

	actions.AddAction("eat", 1, true, goapai.Conditions{
		&goapai.ConditionFn{Key: ATTRIBUTE_3, CheckFn: func(sensors goapai.Sensors) bool {
			if sensors.GetSensor("entity").(*Entity).attributes.hungry {
				return true
			}
			return false
		}},
		&goapai.Condition[int]{Key: ATTRIBUTE_3, Value: 1},
		&goapai.ConditionString{Key: ATTRIBUTE_4, Value: "wipattribute4"},
	}, goapai.Effects{
		goapai.EffectBool{
			Key:      ATTRIBUTE_HUNGRY,
			Value:    false,
			Operator: 0,
		},
		goapai.Effect[int]{
			Key:      ATTRIBUTE_5,
			Value:    1,
			Operator: goapai.EFFECT_ARITHMETIC_ADD,
		},
	})

	//goapai.AddAction("action 82", 82, false, goapai.Conditions{
	//	goapai.ConditionBool{Key: ATTRIBUTE_1, Value: true},
	//}, goapai.Effect{data: map[goapai.StateKey]goapai.StateInterface{ATTRIBUTE_83: goapai.StateBool{Value: false}}}),
	//
	//goapai.AddAction("action 83", 83, false, goapai.Conditions{
	//	goapai.ConditionBool{Key: ATTRIBUTE_1, Value: true},
	//}, goapai.Effect{data: map[goapai.StateKey]goapai.StateInterface{}}),
	//
	//goapai.AddAction("action 84", 84, false, goapai.Conditions{
	//	goapai.ConditionBool{Key: ATTRIBUTE_1, Value: true},
	//}, goapai.Effect{data: map[goapai.StateKey]goapai.StateInterface{ATTRIBUTE_85: goapai.StateBool{Value: false}}}),
	//
	//goapai.AddAction("action 85", 85, false, goapai.Conditions{
	//	goapai.ConditionBool{Key: ATTRIBUTE_1, Value: true},
	//}, goapai.Effect{data: map[goapai.StateKey]goapai.StateInterface{ATTRIBUTE_86: goapai.StateBool{Value: false}}}),
	//
	//goapai.AddAction("action 86", 86, false, goapai.Conditions{
	//	goapai.ConditionBool{Key: ATTRIBUTE_1, Value: true},
	//}, goapai.Effect{data: map[goapai.StateKey]goapai.StateInterface{}}),
	//
	//goapai.AddAction("action 87", 87, false, goapai.Conditions{
	//	goapai.ConditionBool{Key: ATTRIBUTE_1, Value: true},
	//}, goapai.Effect{data: map[goapai.StateKey]goapai.StateInterface{ATTRIBUTE_88: goapai.StateBool{Value: false}}}),
	//
	//goapai.AddAction("action 88", 88, false, goapai.Conditions{
	//	goapai.ConditionBool{Key: ATTRIBUTE_1, Value: true},
	//}, goapai.Effect{data: map[goapai.StateKey]goapai.StateInterface{ATTRIBUTE_89: goapai.StateBool{Value: false}}}),
	//
	//goapai.AddAction("action 89", 89, false, goapai.Conditions{
	//	goapai.ConditionBool{Key: ATTRIBUTE_1, Value: true},
	//}, goapai.Effect{data: map[goapai.StateKey]goapai.StateInterface{}}),
	//
	//goapai.AddAction("action 90", 90, false, goapai.Conditions{
	//	goapai.ConditionBool{Key: ATTRIBUTE_1, Value: true},
	//}, goapai.Effect{data: map[goapai.StateKey]goapai.StateInterface{ATTRIBUTE_90: goapai.StateBool{Value: false}}}),
	//
	//goapai.AddAction("steal an apple", 3, false, goapai.Conditions{
	//	goapai.ConditionBool{Key: ATTRIBUTE_HUNGRY, Value: true},
	//},
	//	goapai.Effect{data: map[goapai.StateKey]goapai.StateInterface{SEE_APPLE: goapai.StateBool{Value: false}, HAVE_APPLE: goapai.StateBool{Value: true}}}),
	//
	//goapai.AddAction("search an apple", 1, false, goapai.Conditions{
	//	goapai.ConditionBool{Key: ATTRIBUTE_HUNGRY, Value: true},
	//},
	//	goapai.Effect{data: map[goapai.StateKey]goapai.StateInterface{SEE_APPLE: goapai.StateBool{Value: true}}}),
	//
	//goapai.AddAction("pick an apple", 1, false, goapai.Conditions{
	//	goapai.ConditionBool{Key: SEE_APPLE, Value: true},
	//}, goapai.Effect{data: map[goapai.StateKey]goapai.StateInterface{SEE_APPLE: goapai.StateBool{Value: false}, HAVE_APPLE: goapai.StateBool{Value: true}}}),
	//
	//goapai.AddAction("eat an apple", 1, false, goapai.Conditions{
	//	goapai.ConditionBool{Key: ATTRIBUTE_HUNGRY, Value: true},
	//	goapai.ConditionBool{Key: HAVE_APPLE, Value: true},
	//}, goapai.Effect{data: map[goapai.StateKey]goapai.StateInterface{ATTRIBUTE_99: goapai.StateBool{Value: true}}}),

	//goapai.AddAction("cook meal", 2.8, false, goapai.Conditions{goapai.ConditionBool{Key: ATTRIBUTE_99, Value: true}}, goapai.Effect{data: map[goapai.StateKey]goapai.StateInterface{HAVE_MEAL: goapai.StateBool{Value: true}}}),
	//goapai.AddAction("eat meal", 1, false, goapai.Conditions{goapai.ConditionBool{Key: HAVE_MEAL, Value: true}}, goapai.Effect{data: map[goapai.StateKey]goapai.StateInterface{ATTRIBUTE_HUNGRY: goapai.StateBool{Value: false}, HAVE_MEAL: goapai.StateBool{Value: false}}}),

	entity := Entity{
		attributes: Attributes{
			happiness:  20,
			hungry:     true,
			hungriness: 80,
			sleepy:     false,
		},
	}

	goals := goapai.Goals{
		"eat": {
			Conditions: []goapai.ConditionInterface{
				&goapai.ConditionBool{Key: ATTRIBUTE_HUNGRY, Value: false},
				&goapai.Condition[int]{Key: ATTRIBUTE_5, Value: 2, Operator: goapai.STATE_OPERATOR_UPPER_OR_EQUAL},
			},
			PriorityFn: func(sensors goapai.Sensors) float64 {
				if sensors.GetSensor("entity").(*Entity).attributes.hungry {
					return 0.8
				}

				return 0.0
			},
		},
		"play": {
			Conditions: []goapai.ConditionInterface{
				&goapai.Condition[int]{Key: ATTRIBUTE_6, Value: 10, Operator: goapai.STATE_OPERATOR_UPPER_OR_EQUAL},
			},
			PriorityFn: func(sensors goapai.Sensors) float64 {
				if sensors.GetSensor("entity").(*Entity).attributes.happiness < 50 {
					return 1.0
				}

				return 0.0
			},
		},
	}
	entity.agent = goapai.CreateAgent(goals, actions)

	goapai.SetSensor(&entity.agent, "entity", &entity)

	goapai.SetStateBool(&entity.agent, ATTRIBUTE_HUNGRY, entity.attributes.hungry)
	goapai.SetStateBool(&entity.agent, ATTRIBUTE_1, entity.attributes.attribute)
	goapai.SetStateBool(&entity.agent, ATTRIBUTE_2, entity.attributes.attribute)
	goapai.SetState(&entity.agent, ATTRIBUTE_3, 1)
	goapai.SetStateString(&entity.agent, ATTRIBUTE_4, "wip")
	//goapai.SetStateBool(entity.agent.states, ATTRIBUTE_5, entity.attributes.attribute)
	goapai.SetState(&entity.agent, ATTRIBUTE_6, 0)
	goapai.SetStateBool(&entity.agent, ATTRIBUTE_7, entity.attributes.attribute)
	goapai.SetStateBool(&entity.agent, ATTRIBUTE_8, entity.attributes.attribute)
	goapai.SetStateBool(&entity.agent, ATTRIBUTE_9, entity.attributes.attribute)
	goapai.SetStateBool(&entity.agent, ATTRIBUTE_10, entity.attributes.attribute)
	goapai.SetStateBool(&entity.agent, ATTRIBUTE_11, entity.attributes.attribute)
	goapai.SetStateBool(&entity.agent, ATTRIBUTE_12, entity.attributes.attribute)
	goapai.SetStateBool(&entity.agent, ATTRIBUTE_13, entity.attributes.attribute)
	goapai.SetStateBool(&entity.agent, ATTRIBUTE_14, entity.attributes.attribute)
	goapai.SetStateBool(&entity.agent, ATTRIBUTE_15, entity.attributes.attribute)
	goapai.SetStateBool(&entity.agent, ATTRIBUTE_16, entity.attributes.attribute)
	goapai.SetStateBool(&entity.agent, ATTRIBUTE_17, entity.attributes.attribute)
	goapai.SetStateBool(&entity.agent, ATTRIBUTE_18, entity.attributes.attribute)
	goapai.SetStateBool(&entity.agent, ATTRIBUTE_19, entity.attributes.attribute)
	goapai.SetStateBool(&entity.agent, ATTRIBUTE_20, entity.attributes.attribute)
	goapai.SetStateBool(&entity.agent, ATTRIBUTE_21, entity.attributes.attribute)
	goapai.SetStateBool(&entity.agent, ATTRIBUTE_22, entity.attributes.attribute)
	goapai.SetStateBool(&entity.agent, ATTRIBUTE_23, entity.attributes.attribute)
	goapai.SetStateBool(&entity.agent, ATTRIBUTE_24, entity.attributes.attribute)
	goapai.SetStateBool(&entity.agent, ATTRIBUTE_25, entity.attributes.attribute)
	goapai.SetStateBool(&entity.agent, ATTRIBUTE_26, entity.attributes.attribute)
	goapai.SetStateBool(&entity.agent, ATTRIBUTE_27, entity.attributes.attribute)
	goapai.SetStateBool(&entity.agent, ATTRIBUTE_28, entity.attributes.attribute)
	goapai.SetStateBool(&entity.agent, ATTRIBUTE_29, entity.attributes.attribute)
	goapai.SetStateBool(&entity.agent, ATTRIBUTE_30, entity.attributes.attribute)
	goapai.SetStateBool(&entity.agent, ATTRIBUTE_31, entity.attributes.attribute)
	goapai.SetStateBool(&entity.agent, ATTRIBUTE_32, entity.attributes.attribute)
	goapai.SetStateBool(&entity.agent, ATTRIBUTE_33, entity.attributes.attribute)
	goapai.SetStateBool(&entity.agent, ATTRIBUTE_34, entity.attributes.attribute)
	goapai.SetStateBool(&entity.agent, ATTRIBUTE_35, entity.attributes.attribute)
	goapai.SetStateBool(&entity.agent, ATTRIBUTE_36, entity.attributes.attribute)
	goapai.SetStateBool(&entity.agent, ATTRIBUTE_37, entity.attributes.attribute)
	goapai.SetStateBool(&entity.agent, ATTRIBUTE_38, entity.attributes.attribute)
	goapai.SetStateBool(&entity.agent, ATTRIBUTE_39, entity.attributes.attribute)
	goapai.SetStateBool(&entity.agent, ATTRIBUTE_40, entity.attributes.attribute)
	goapai.SetStateBool(&entity.agent, ATTRIBUTE_41, entity.attributes.attribute)
	goapai.SetStateBool(&entity.agent, ATTRIBUTE_42, entity.attributes.attribute)
	goapai.SetStateBool(&entity.agent, ATTRIBUTE_43, entity.attributes.attribute)
	goapai.SetStateBool(&entity.agent, ATTRIBUTE_44, entity.attributes.attribute)
	goapai.SetStateBool(&entity.agent, ATTRIBUTE_45, entity.attributes.attribute)
	goapai.SetStateBool(&entity.agent, ATTRIBUTE_46, entity.attributes.attribute)
	goapai.SetStateBool(&entity.agent, ATTRIBUTE_47, entity.attributes.attribute)
	goapai.SetStateBool(&entity.agent, ATTRIBUTE_48, entity.attributes.attribute)
	goapai.SetStateBool(&entity.agent, ATTRIBUTE_49, entity.attributes.attribute)
	goapai.SetStateBool(&entity.agent, ATTRIBUTE_50, entity.attributes.attribute)
	goapai.SetStateBool(&entity.agent, ATTRIBUTE_51, entity.attributes.attribute)
	goapai.SetStateBool(&entity.agent, ATTRIBUTE_52, entity.attributes.attribute)
	goapai.SetStateBool(&entity.agent, ATTRIBUTE_53, entity.attributes.attribute)
	goapai.SetStateBool(&entity.agent, ATTRIBUTE_54, entity.attributes.attribute)
	goapai.SetStateBool(&entity.agent, ATTRIBUTE_55, entity.attributes.attribute)
	goapai.SetStateBool(&entity.agent, ATTRIBUTE_56, entity.attributes.attribute)
	goapai.SetStateBool(&entity.agent, ATTRIBUTE_57, entity.attributes.attribute)
	goapai.SetStateBool(&entity.agent, ATTRIBUTE_58, entity.attributes.attribute)
	goapai.SetStateBool(&entity.agent, ATTRIBUTE_59, entity.attributes.attribute)
	goapai.SetStateBool(&entity.agent, ATTRIBUTE_60, entity.attributes.attribute)
	goapai.SetStateBool(&entity.agent, ATTRIBUTE_61, entity.attributes.attribute)
	goapai.SetStateBool(&entity.agent, ATTRIBUTE_62, entity.attributes.attribute)
	goapai.SetStateBool(&entity.agent, ATTRIBUTE_63, entity.attributes.attribute)
	goapai.SetStateBool(&entity.agent, ATTRIBUTE_64, entity.attributes.attribute)
	goapai.SetStateBool(&entity.agent, ATTRIBUTE_65, entity.attributes.attribute)
	goapai.SetStateBool(&entity.agent, ATTRIBUTE_66, entity.attributes.attribute)
	goapai.SetStateBool(&entity.agent, ATTRIBUTE_67, entity.attributes.attribute)
	goapai.SetStateBool(&entity.agent, ATTRIBUTE_68, entity.attributes.attribute)
	goapai.SetStateBool(&entity.agent, ATTRIBUTE_69, entity.attributes.attribute)
	goapai.SetStateBool(&entity.agent, ATTRIBUTE_70, entity.attributes.attribute)
	goapai.SetStateBool(&entity.agent, ATTRIBUTE_71, entity.attributes.attribute)
	goapai.SetStateBool(&entity.agent, ATTRIBUTE_72, entity.attributes.attribute)
	goapai.SetStateBool(&entity.agent, ATTRIBUTE_73, entity.attributes.attribute)
	goapai.SetStateBool(&entity.agent, ATTRIBUTE_74, entity.attributes.attribute)
	goapai.SetStateBool(&entity.agent, ATTRIBUTE_75, entity.attributes.attribute)
	goapai.SetStateBool(&entity.agent, ATTRIBUTE_76, entity.attributes.attribute)
	goapai.SetStateBool(&entity.agent, ATTRIBUTE_77, entity.attributes.attribute)
	goapai.SetStateBool(&entity.agent, ATTRIBUTE_78, entity.attributes.attribute)
	goapai.SetStateBool(&entity.agent, ATTRIBUTE_79, entity.attributes.attribute)
	goapai.SetStateBool(&entity.agent, ATTRIBUTE_80, entity.attributes.attribute)
	goapai.SetStateBool(&entity.agent, ATTRIBUTE_81, entity.attributes.attribute)
	goapai.SetStateBool(&entity.agent, ATTRIBUTE_82, entity.attributes.attribute)
	goapai.SetStateBool(&entity.agent, ATTRIBUTE_83, entity.attributes.attribute)
	goapai.SetStateBool(&entity.agent, ATTRIBUTE_84, entity.attributes.attribute)
	goapai.SetStateBool(&entity.agent, ATTRIBUTE_85, entity.attributes.attribute)
	goapai.SetStateBool(&entity.agent, ATTRIBUTE_86, entity.attributes.attribute)
	goapai.SetStateBool(&entity.agent, ATTRIBUTE_87, entity.attributes.attribute)
	goapai.SetStateBool(&entity.agent, ATTRIBUTE_88, entity.attributes.attribute)
	goapai.SetStateBool(&entity.agent, ATTRIBUTE_89, entity.attributes.attribute)
	goapai.SetStateBool(&entity.agent, ATTRIBUTE_90, entity.attributes.attribute)
	goapai.SetStateBool(&entity.agent, ATTRIBUTE_91, entity.attributes.attribute)
	goapai.SetStateBool(&entity.agent, ATTRIBUTE_92, entity.attributes.attribute)
	goapai.SetStateBool(&entity.agent, ATTRIBUTE_93, entity.attributes.attribute)
	goapai.SetStateBool(&entity.agent, ATTRIBUTE_94, entity.attributes.attribute)
	goapai.SetStateBool(&entity.agent, ATTRIBUTE_95, entity.attributes.attribute)
	goapai.SetStateBool(&entity.agent, ATTRIBUTE_96, entity.attributes.attribute)
	goapai.SetStateBool(&entity.agent, ATTRIBUTE_97, entity.attributes.attribute)
	goapai.SetStateBool(&entity.agent, ATTRIBUTE_98, entity.attributes.attribute)
	goapai.SetStateBool(&entity.agent, ATTRIBUTE_99, entity.attributes.attribute)
	goapai.SetStateBool(&entity.agent, ATTRIBUTE_100, entity.attributes.attribute)

	// Write to the trace file.
	f, _ := os.Create("trace.out")
	fcpu, _ := os.Create(`cpu.prof`)
	fheap, _ := os.Create(`heap.prof`)

	pprof.StartCPUProfile(fcpu)
	pprof.WriteHeapProfile(fheap)
	trace.Start(f)

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_, _ = goapai.GetPlan(entity.agent, 10)
	}

	t := time.Now()
	_, plan := goapai.GetPlan(entity.agent, 10)
	for _, a := range plan {
		fmt.Println(a.GetName())
	}
	fmt.Println("Total cost:", plan.GetTotalCost())
	fmt.Println("elapsed time", time.Now().Sub(t))

	defer f.Close()
	defer fcpu.Close()
	defer fheap.Close()

	trace.Stop()
	pprof.StopCPUProfile()

	b.ReportAllocs()
}
