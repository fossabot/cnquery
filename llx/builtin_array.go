package llx

import (
	"errors"
	"strconv"

	"go.mondoo.io/mondoo/types"
)

var arrayBlockType = types.Array(types.Map(types.String, types.Any))

// arrayFunctions are all the handlers for builtin array methods
var arrayFunctions map[string]chunkHandler

func arrayGetIndex(c *LeiseExecutor, bind *RawData, chunk *Chunk, ref int32) (*RawData, int32, error) {
	if bind.Value == nil {
		return &RawData{Type: bind.Type[1:]}, 0, nil
	}

	args := chunk.Function.Args
	// TODO: all this needs to go into the compile phase
	if len(args) < 1 {
		return nil, 0, errors.New("Called [] with " + strconv.Itoa(len(args)) + " arguments, only 1 supported.")
	}
	if len(args) > 1 {
		return nil, 0, errors.New("called [] with " + strconv.Itoa(len(args)) + " arguments, only 1 supported.")
	}
	t := types.Type(args[0].Type)
	if t != types.Int {
		return nil, 0, errors.New("called [] with wrong type " + t.Label())
	}
	// ^^ TODO

	key := int(bytes2int(args[0].Value))

	arr, ok := bind.Value.([]interface{})
	if !ok {
		return nil, 0, errors.New("failed to typecast " + bind.Type.Label() + " into array")
	}

	if key < 0 {
		return nil, 0, errors.New("array index out of bound (trying to access element " + strconv.Itoa(key) + ")")
	}
	if key >= len(arr) {
		return nil, 0, errors.New("array index out of bound (trying to access element " + strconv.Itoa(key) + ", max: " + strconv.Itoa(len(arr)-1) + ")")
	}

	return &RawData{
		Type:  bind.Type[1:],
		Value: arr[key],
	}, 0, nil
}

func arrayBlockList(c *LeiseExecutor, bind *RawData, chunk *Chunk, ref int32) (*RawData, int32, error) {
	arr, ok := bind.Value.([]interface{})
	if !ok {
		return nil, 0, errors.New("failed to typecast " + bind.Type.Label() + " into array")
	}

	if len(arr) == 0 {
		return bind, 0, nil
	}

	prim := chunk.Function.Args[0]
	if !types.Type(prim.Type).IsFunction() {
		return nil, 0, errors.New("called block with wrong function type")
	}
	fref, ok := prim.Ref()
	if !ok {
		return nil, 0, errors.New("cannot retrieve function reference on block call")
	}
	fun := c.code.Functions[fref-1]
	if fun == nil {
		return nil, 0, errors.New("block function is nil")
	}

	// pre-init everything to avoid concurrency issues with long list
	allResults := make([]interface{}, len(arr))
	for idx := range arr {
		blockResult := map[string]interface{}{}
		allResults[idx] = blockResult
	}

	var childType types.Type
	if bind.Type == types.Dict {
		childType = types.Dict
	} else {
		childType = bind.Type.Child()
	}

	finishedBlocks := 0

	for idx := range arr {
		blockResult := allResults[idx].(map[string]interface{})

		bind := &RawData{
			Type:  childType,
			Value: arr[idx],
		}

		finished := false
		err := c.runFunctionBlock(bind, fun, func(res *RawResult) {
			blockResult[res.CodeID] = res.Data

			if len(blockResult) == len(fun.Entrypoints) && !finished {
				finishedBlocks++
				finished = true
			}

			if finishedBlocks >= len(arr) {
				c.cache.Store(ref, &stepCache{
					Result: &RawData{
						Type:  arrayBlockType,
						Value: allResults,
					},
					IsStatic: true,
				})
				c.triggerChain(ref)
			}
		})
		if err != nil {
			return nil, 0, err
		}
	}

	return nil, 0, nil
}

func arrayWhere(c *LeiseExecutor, bind *RawData, chunk *Chunk, ref int32) (*RawData, int32, error) {
	// where(array, function)
	itemsRef := chunk.Function.Args[0]
	items, rref, err := c.resolveValue(itemsRef, ref)
	if err != nil || rref > 0 {
		return nil, rref, err
	}

	if items.Value == nil {
		return &RawData{Type: items.Type}, 0, nil
	}

	list := items.Value.([]interface{})
	if len(list) == 0 {
		return items, 0, nil
	}

	arg1 := chunk.Function.Args[1]
	fref, ok := arg1.Ref()
	if !ok {
		return nil, 0, errors.New("Failed to retrieve function reference of 'where' call")
	}

	f := c.code.Functions[fref-1]
	ct := items.Type.Child()
	filteredList := map[int]interface{}{}
	finishedResults := 0
	for i := range list {
		c.runFunctionBlock(&RawData{Type: ct, Value: list[i]}, f, func(res *RawResult) {
			_, ok := filteredList[i]
			if !ok {
				finishedResults++
			}

			isTruthy, _ := res.Data.IsTruthy()
			if isTruthy {
				filteredList[i] = list[i]
			} else {
				filteredList[i] = nil
			}

			if finishedResults == len(list) {
				resList := []interface{}{}
				for j := 0; j < len(filteredList); j++ {
					k := filteredList[j]
					if k != nil {
						resList = append(resList, k)
					}
				}

				c.cache.Store(ref, &stepCache{
					Result: &RawData{
						Type:  bind.Type,
						Value: resList,
					},
					IsStatic: false,
				})
				c.triggerChain(ref)
			}
		})
	}

	return nil, 0, nil
}

func arrayLength(c *LeiseExecutor, bind *RawData, chunk *Chunk, ref int32) (*RawData, int32, error) {
	if bind.Value == nil {
		return &RawData{Type: types.Int}, 0, nil
	}

	arr, ok := bind.Value.([]interface{})
	if !ok {
		return nil, 0, errors.New("failed to typecast " + bind.Type.Label() + " into array")
	}
	return IntData(int64(len(arr))), 0, nil
}

func compileArrayOpArray(op string) func(types.Type, types.Type) (string, error) {
	return func(left types.Type, right types.Type) (string, error) {
		name := string(left.Child()) + op + string(right)
		af := BuiltinFunctions[types.ArrayLike]
		if _, ok := af[name]; ok {
			return name, nil
		}

		return "<T>" + op + "<T>", nil
	}
}

func compileLogicalArrayOp(underlying types.Type, op string) func(types.Type, types.Type) (string, error) {
	return func(left types.Type, right types.Type) (string, error) {
		name := string(types.Any) + op + string(right.Underlying())
		af := BuiltinFunctions[underlying]
		if _, ok := af[name]; ok {
			return name, nil
		}

		return "", errors.New("cannot find operation for " + left.Label() + " " + op + " " + right.Label())
	}
}

func cmpArrays(left *RawData, right *RawData, f func(interface{}, interface{}) bool) bool {
	if left.Value == nil {
		if right.Value == nil {
			return true
		}
		return false
	}
	if right.Value == nil {
		return false
	}

	l := left.Value.([]interface{})
	r := right.Value.([]interface{})

	if len(l) != len(r) {
		return false
	}

	for i := range l {
		if !f(l[i], r[i]) {
			return false
		}
	}

	return true
}

func cmpArrayOne(leftArray *RawData, right *RawData, f func(interface{}, interface{}) bool) bool {
	l := leftArray.Value.([]interface{})
	if len(l) != 1 {
		return false
	}
	return f(l[0], right.Value)
}

// []T -- []T

func tArrayCmp(left *RawData, right *RawData) func(interface{}, interface{}) bool {
	return func(a interface{}, b interface{}) bool {
		if left.Type.Child() != right.Type.Child() {
			return false
		}
		return a == b
	}
}

func tarrayCmpTarray(c *LeiseExecutor, bind *RawData, chunk *Chunk, ref int32) (*RawData, int32, error) {
	return rawboolOp(c, bind, chunk, ref, func(left *RawData, right *RawData) bool {
		return cmpArrays(left, right, tArrayCmp(left, right))
	})
}

func tarrayNotTarray(c *LeiseExecutor, bind *RawData, chunk *Chunk, ref int32) (*RawData, int32, error) {
	return rawboolNotOp(c, bind, chunk, ref, func(left *RawData, right *RawData) bool {
		return cmpArrays(left, right, tArrayCmp(left, right))
	})
}

func boolarrayCmpBoolarray(c *LeiseExecutor, bind *RawData, chunk *Chunk, ref int32) (*RawData, int32, error) {
	return rawboolOp(c, bind, chunk, ref, func(left *RawData, right *RawData) bool {
		return cmpArrays(left, right, opBoolCmpBool)
	})
}

func intarrayCmpIntarray(c *LeiseExecutor, bind *RawData, chunk *Chunk, ref int32) (*RawData, int32, error) {
	return rawboolOp(c, bind, chunk, ref, func(left *RawData, right *RawData) bool {
		return cmpArrays(left, right, opIntCmpInt)
	})
}

func floatarrayCmpFloatarray(c *LeiseExecutor, bind *RawData, chunk *Chunk, ref int32) (*RawData, int32, error) {
	return rawboolOp(c, bind, chunk, ref, func(left *RawData, right *RawData) bool {
		return cmpArrays(left, right, opFloatCmpFloat)
	})
}

func stringarrayCmpStringarray(c *LeiseExecutor, bind *RawData, chunk *Chunk, ref int32) (*RawData, int32, error) {
	return rawboolOp(c, bind, chunk, ref, func(left *RawData, right *RawData) bool {
		return cmpArrays(left, right, opStringCmpString)
	})
}

func boolarrayNotBoolarray(c *LeiseExecutor, bind *RawData, chunk *Chunk, ref int32) (*RawData, int32, error) {
	return rawboolNotOp(c, bind, chunk, ref, func(left *RawData, right *RawData) bool {
		return cmpArrays(left, right, opBoolCmpBool)
	})
}

func intarrayNotIntarray(c *LeiseExecutor, bind *RawData, chunk *Chunk, ref int32) (*RawData, int32, error) {
	return rawboolNotOp(c, bind, chunk, ref, func(left *RawData, right *RawData) bool {
		return cmpArrays(left, right, opIntCmpInt)
	})
}

func floatarrayNotFloatarray(c *LeiseExecutor, bind *RawData, chunk *Chunk, ref int32) (*RawData, int32, error) {
	return rawboolNotOp(c, bind, chunk, ref, func(left *RawData, right *RawData) bool {
		return cmpArrays(left, right, opFloatCmpFloat)
	})
}

func stringarrayNotStringarray(c *LeiseExecutor, bind *RawData, chunk *Chunk, ref int32) (*RawData, int32, error) {
	return rawboolNotOp(c, bind, chunk, ref, func(left *RawData, right *RawData) bool {
		return cmpArrays(left, right, opStringCmpString)
	})
}

// []T -- T

func boolarrayCmpBool(c *LeiseExecutor, bind *RawData, chunk *Chunk, ref int32) (*RawData, int32, error) {
	return rawboolOp(c, bind, chunk, ref, func(left *RawData, right *RawData) bool {
		return cmpArrayOne(left, right, opBoolCmpBool)
	})
}

func boolarrayNotBool(c *LeiseExecutor, bind *RawData, chunk *Chunk, ref int32) (*RawData, int32, error) {
	return rawboolNotOp(c, bind, chunk, ref, func(left *RawData, right *RawData) bool {
		return cmpArrayOne(left, right, opBoolCmpBool)
	})
}

func intarrayCmpInt(c *LeiseExecutor, bind *RawData, chunk *Chunk, ref int32) (*RawData, int32, error) {
	return rawboolOp(c, bind, chunk, ref, func(left *RawData, right *RawData) bool {
		return cmpArrayOne(left, right, opIntCmpInt)
	})
}

func intarrayNotInt(c *LeiseExecutor, bind *RawData, chunk *Chunk, ref int32) (*RawData, int32, error) {
	return rawboolNotOp(c, bind, chunk, ref, func(left *RawData, right *RawData) bool {
		return cmpArrayOne(left, right, opIntCmpInt)
	})
}

func floatarrayCmpFloat(c *LeiseExecutor, bind *RawData, chunk *Chunk, ref int32) (*RawData, int32, error) {
	return rawboolOp(c, bind, chunk, ref, func(left *RawData, right *RawData) bool {
		return cmpArrayOne(left, right, opFloatCmpFloat)
	})
}

func floatarrayNotFloat(c *LeiseExecutor, bind *RawData, chunk *Chunk, ref int32) (*RawData, int32, error) {
	return rawboolNotOp(c, bind, chunk, ref, func(left *RawData, right *RawData) bool {
		return cmpArrayOne(left, right, opFloatCmpFloat)
	})
}

func stringarrayCmpString(c *LeiseExecutor, bind *RawData, chunk *Chunk, ref int32) (*RawData, int32, error) {
	return rawboolOp(c, bind, chunk, ref, func(left *RawData, right *RawData) bool {
		return cmpArrayOne(left, right, opStringCmpString)
	})
}

func stringarrayNotString(c *LeiseExecutor, bind *RawData, chunk *Chunk, ref int32) (*RawData, int32, error) {
	return rawboolNotOp(c, bind, chunk, ref, func(left *RawData, right *RawData) bool {
		return cmpArrayOne(left, right, opStringCmpString)
	})
}

// T -- []T

func boolCmpBoolarray(c *LeiseExecutor, bind *RawData, chunk *Chunk, ref int32) (*RawData, int32, error) {
	return rawboolOp(c, bind, chunk, ref, func(left *RawData, right *RawData) bool {
		return cmpArrayOne(right, left, opBoolCmpBool)
	})
}

func boolNotBoolarray(c *LeiseExecutor, bind *RawData, chunk *Chunk, ref int32) (*RawData, int32, error) {
	return rawboolNotOp(c, bind, chunk, ref, func(left *RawData, right *RawData) bool {
		return cmpArrayOne(right, left, opBoolCmpBool)
	})
}

func intCmpIntarray(c *LeiseExecutor, bind *RawData, chunk *Chunk, ref int32) (*RawData, int32, error) {
	return rawboolOp(c, bind, chunk, ref, func(left *RawData, right *RawData) bool {
		return cmpArrayOne(right, left, opIntCmpInt)
	})
}

func intNotIntarray(c *LeiseExecutor, bind *RawData, chunk *Chunk, ref int32) (*RawData, int32, error) {
	return rawboolNotOp(c, bind, chunk, ref, func(left *RawData, right *RawData) bool {
		return cmpArrayOne(right, left, opIntCmpInt)
	})
}

func floatCmpFloatarray(c *LeiseExecutor, bind *RawData, chunk *Chunk, ref int32) (*RawData, int32, error) {
	return rawboolOp(c, bind, chunk, ref, func(left *RawData, right *RawData) bool {
		return cmpArrayOne(right, left, opFloatCmpFloat)
	})
}

func floatNotFloatarray(c *LeiseExecutor, bind *RawData, chunk *Chunk, ref int32) (*RawData, int32, error) {
	return rawboolNotOp(c, bind, chunk, ref, func(left *RawData, right *RawData) bool {
		return cmpArrayOne(right, left, opFloatCmpFloat)
	})
}

func stringCmpStringarray(c *LeiseExecutor, bind *RawData, chunk *Chunk, ref int32) (*RawData, int32, error) {
	return rawboolOp(c, bind, chunk, ref, func(left *RawData, right *RawData) bool {
		return cmpArrayOne(right, left, opStringCmpString)
	})
}

func stringNotStringarray(c *LeiseExecutor, bind *RawData, chunk *Chunk, ref int32) (*RawData, int32, error) {
	return rawboolNotOp(c, bind, chunk, ref, func(left *RawData, right *RawData) bool {
		return cmpArrayOne(right, left, opStringCmpString)
	})
}

// string -- []T

func stringCmpBoolarray(c *LeiseExecutor, bind *RawData, chunk *Chunk, ref int32) (*RawData, int32, error) {
	return rawboolOp(c, bind, chunk, ref, func(left *RawData, right *RawData) bool {
		return cmpArrayOne(right, left, opBoolCmpString)
	})
}

func stringNotBoolarray(c *LeiseExecutor, bind *RawData, chunk *Chunk, ref int32) (*RawData, int32, error) {
	return rawboolNotOp(c, bind, chunk, ref, func(left *RawData, right *RawData) bool {
		return cmpArrayOne(right, left, opBoolCmpString)
	})
}

func boolarrayCmpString(c *LeiseExecutor, bind *RawData, chunk *Chunk, ref int32) (*RawData, int32, error) {
	return rawboolOp(c, bind, chunk, ref, func(left *RawData, right *RawData) bool {
		return cmpArrayOne(left, right, opBoolCmpString)
	})
}

func boolarrayNotString(c *LeiseExecutor, bind *RawData, chunk *Chunk, ref int32) (*RawData, int32, error) {
	return rawboolNotOp(c, bind, chunk, ref, func(left *RawData, right *RawData) bool {
		return cmpArrayOne(left, right, opBoolCmpString)
	})
}

func stringCmpIntarray(c *LeiseExecutor, bind *RawData, chunk *Chunk, ref int32) (*RawData, int32, error) {
	return rawboolOp(c, bind, chunk, ref, func(left *RawData, right *RawData) bool {
		return cmpArrayOne(right, left, opIntCmpString)
	})
}

func stringNotIntarray(c *LeiseExecutor, bind *RawData, chunk *Chunk, ref int32) (*RawData, int32, error) {
	return rawboolNotOp(c, bind, chunk, ref, func(left *RawData, right *RawData) bool {
		return cmpArrayOne(right, left, opIntCmpString)
	})
}

func intarrayCmpString(c *LeiseExecutor, bind *RawData, chunk *Chunk, ref int32) (*RawData, int32, error) {
	return rawboolOp(c, bind, chunk, ref, func(left *RawData, right *RawData) bool {
		return cmpArrayOne(left, right, opIntCmpString)
	})
}

func intarrayNotString(c *LeiseExecutor, bind *RawData, chunk *Chunk, ref int32) (*RawData, int32, error) {
	return rawboolNotOp(c, bind, chunk, ref, func(left *RawData, right *RawData) bool {
		return cmpArrayOne(left, right, opIntCmpString)
	})
}

func stringCmpFloatarray(c *LeiseExecutor, bind *RawData, chunk *Chunk, ref int32) (*RawData, int32, error) {
	return rawboolOp(c, bind, chunk, ref, func(left *RawData, right *RawData) bool {
		return cmpArrayOne(right, left, opFloatCmpString)
	})
}

func stringNotFloatarray(c *LeiseExecutor, bind *RawData, chunk *Chunk, ref int32) (*RawData, int32, error) {
	return rawboolNotOp(c, bind, chunk, ref, func(left *RawData, right *RawData) bool {
		return cmpArrayOne(right, left, opFloatCmpString)
	})
}

func floatarrayCmpString(c *LeiseExecutor, bind *RawData, chunk *Chunk, ref int32) (*RawData, int32, error) {
	return rawboolOp(c, bind, chunk, ref, func(left *RawData, right *RawData) bool {
		return cmpArrayOne(left, right, opFloatCmpString)
	})
}

func floatarrayNotString(c *LeiseExecutor, bind *RawData, chunk *Chunk, ref int32) (*RawData, int32, error) {
	return rawboolNotOp(c, bind, chunk, ref, func(left *RawData, right *RawData) bool {
		return cmpArrayOne(left, right, opFloatCmpString)
	})
}

// bool -- []string

func boolCmpStringarray(c *LeiseExecutor, bind *RawData, chunk *Chunk, ref int32) (*RawData, int32, error) {
	return rawboolOp(c, bind, chunk, ref, func(left *RawData, right *RawData) bool {
		return cmpArrayOne(right, left, opStringCmpBool)
	})
}

func boolNotStringarray(c *LeiseExecutor, bind *RawData, chunk *Chunk, ref int32) (*RawData, int32, error) {
	return rawboolNotOp(c, bind, chunk, ref, func(left *RawData, right *RawData) bool {
		return cmpArrayOne(right, left, opStringCmpBool)
	})
}

func stringarrayCmpBool(c *LeiseExecutor, bind *RawData, chunk *Chunk, ref int32) (*RawData, int32, error) {
	return rawboolOp(c, bind, chunk, ref, func(left *RawData, right *RawData) bool {
		return cmpArrayOne(left, right, opStringCmpBool)
	})
}

func stringarrayNotBool(c *LeiseExecutor, bind *RawData, chunk *Chunk, ref int32) (*RawData, int32, error) {
	return rawboolNotOp(c, bind, chunk, ref, func(left *RawData, right *RawData) bool {
		return cmpArrayOne(left, right, opStringCmpBool)
	})
}

// int -- []string

func intCmpStringarray(c *LeiseExecutor, bind *RawData, chunk *Chunk, ref int32) (*RawData, int32, error) {
	return rawboolOp(c, bind, chunk, ref, func(left *RawData, right *RawData) bool {
		return cmpArrayOne(right, left, opStringCmpInt)
	})
}

func intNotStringarray(c *LeiseExecutor, bind *RawData, chunk *Chunk, ref int32) (*RawData, int32, error) {
	return rawboolNotOp(c, bind, chunk, ref, func(left *RawData, right *RawData) bool {
		return cmpArrayOne(right, left, opStringCmpInt)
	})
}

func stringarrayCmpInt(c *LeiseExecutor, bind *RawData, chunk *Chunk, ref int32) (*RawData, int32, error) {
	return rawboolOp(c, bind, chunk, ref, func(left *RawData, right *RawData) bool {
		return cmpArrayOne(left, right, opStringCmpInt)
	})
}

func stringarrayNotInt(c *LeiseExecutor, bind *RawData, chunk *Chunk, ref int32) (*RawData, int32, error) {
	return rawboolNotOp(c, bind, chunk, ref, func(left *RawData, right *RawData) bool {
		return cmpArrayOne(left, right, opStringCmpInt)
	})
}

// float -- []string

func floatCmpStringarray(c *LeiseExecutor, bind *RawData, chunk *Chunk, ref int32) (*RawData, int32, error) {
	return rawboolOp(c, bind, chunk, ref, func(left *RawData, right *RawData) bool {
		return cmpArrayOne(right, left, opStringCmpFloat)
	})
}

func floatNotStringarray(c *LeiseExecutor, bind *RawData, chunk *Chunk, ref int32) (*RawData, int32, error) {
	return rawboolNotOp(c, bind, chunk, ref, func(left *RawData, right *RawData) bool {
		return cmpArrayOne(right, left, opStringCmpFloat)
	})
}

func stringarrayCmpFloat(c *LeiseExecutor, bind *RawData, chunk *Chunk, ref int32) (*RawData, int32, error) {
	return rawboolOp(c, bind, chunk, ref, func(left *RawData, right *RawData) bool {
		return cmpArrayOne(left, right, opStringCmpFloat)
	})
}

func stringarrayNotFloat(c *LeiseExecutor, bind *RawData, chunk *Chunk, ref int32) (*RawData, int32, error) {
	return rawboolNotOp(c, bind, chunk, ref, func(left *RawData, right *RawData) bool {
		return cmpArrayOne(left, right, opStringCmpFloat)
	})
}

// regex -- []T

func regexCmpStringarray(c *LeiseExecutor, bind *RawData, chunk *Chunk, ref int32) (*RawData, int32, error) {
	return rawboolOp(c, bind, chunk, ref, func(left *RawData, right *RawData) bool {
		return cmpArrayOne(right, left, opStringCmpRegex)
	})
}

func regexNotStringarray(c *LeiseExecutor, bind *RawData, chunk *Chunk, ref int32) (*RawData, int32, error) {
	return rawboolNotOp(c, bind, chunk, ref, func(left *RawData, right *RawData) bool {
		return cmpArrayOne(right, left, opStringCmpRegex)
	})
}

func stringarrayCmpRegex(c *LeiseExecutor, bind *RawData, chunk *Chunk, ref int32) (*RawData, int32, error) {
	return rawboolOp(c, bind, chunk, ref, func(left *RawData, right *RawData) bool {
		return cmpArrayOne(left, right, opStringCmpRegex)
	})
}

func stringarrayNotRegex(c *LeiseExecutor, bind *RawData, chunk *Chunk, ref int32) (*RawData, int32, error) {
	return rawboolNotOp(c, bind, chunk, ref, func(left *RawData, right *RawData) bool {
		return cmpArrayOne(left, right, opStringCmpRegex)
	})
}

func regexCmpBoolarray(c *LeiseExecutor, bind *RawData, chunk *Chunk, ref int32) (*RawData, int32, error) {
	return rawboolOp(c, bind, chunk, ref, func(left *RawData, right *RawData) bool {
		return cmpArrayOne(right, left, opBoolCmpRegex)
	})
}

func regexNotBoolarray(c *LeiseExecutor, bind *RawData, chunk *Chunk, ref int32) (*RawData, int32, error) {
	return rawboolNotOp(c, bind, chunk, ref, func(left *RawData, right *RawData) bool {
		return cmpArrayOne(right, left, opBoolCmpRegex)
	})
}

func boolarrayCmpRegex(c *LeiseExecutor, bind *RawData, chunk *Chunk, ref int32) (*RawData, int32, error) {
	return rawboolOp(c, bind, chunk, ref, func(left *RawData, right *RawData) bool {
		return cmpArrayOne(left, right, opBoolCmpRegex)
	})
}

func boolarrayNotRegex(c *LeiseExecutor, bind *RawData, chunk *Chunk, ref int32) (*RawData, int32, error) {
	return rawboolNotOp(c, bind, chunk, ref, func(left *RawData, right *RawData) bool {
		return cmpArrayOne(left, right, opBoolCmpRegex)
	})
}

func regexCmpIntarray(c *LeiseExecutor, bind *RawData, chunk *Chunk, ref int32) (*RawData, int32, error) {
	return rawboolOp(c, bind, chunk, ref, func(left *RawData, right *RawData) bool {
		return cmpArrayOne(right, left, opIntCmpRegex)
	})
}

func regexNotIntarray(c *LeiseExecutor, bind *RawData, chunk *Chunk, ref int32) (*RawData, int32, error) {
	return rawboolNotOp(c, bind, chunk, ref, func(left *RawData, right *RawData) bool {
		return cmpArrayOne(right, left, opIntCmpRegex)
	})
}

func intarrayCmpRegex(c *LeiseExecutor, bind *RawData, chunk *Chunk, ref int32) (*RawData, int32, error) {
	return rawboolOp(c, bind, chunk, ref, func(left *RawData, right *RawData) bool {
		return cmpArrayOne(left, right, opIntCmpRegex)
	})
}

func intarrayNotRegex(c *LeiseExecutor, bind *RawData, chunk *Chunk, ref int32) (*RawData, int32, error) {
	return rawboolNotOp(c, bind, chunk, ref, func(left *RawData, right *RawData) bool {
		return cmpArrayOne(left, right, opIntCmpRegex)
	})
}

func regexCmpFloatarray(c *LeiseExecutor, bind *RawData, chunk *Chunk, ref int32) (*RawData, int32, error) {
	return rawboolOp(c, bind, chunk, ref, func(left *RawData, right *RawData) bool {
		return cmpArrayOne(right, left, opFloatCmpRegex)
	})
}

func regexNotFloatarray(c *LeiseExecutor, bind *RawData, chunk *Chunk, ref int32) (*RawData, int32, error) {
	return rawboolNotOp(c, bind, chunk, ref, func(left *RawData, right *RawData) bool {
		return cmpArrayOne(right, left, opFloatCmpRegex)
	})
}

func floatarrayCmpRegex(c *LeiseExecutor, bind *RawData, chunk *Chunk, ref int32) (*RawData, int32, error) {
	return rawboolOp(c, bind, chunk, ref, func(left *RawData, right *RawData) bool {
		return cmpArrayOne(left, right, opFloatCmpRegex)
	})
}

func floatarrayNotRegex(c *LeiseExecutor, bind *RawData, chunk *Chunk, ref int32) (*RawData, int32, error) {
	return rawboolNotOp(c, bind, chunk, ref, func(left *RawData, right *RawData) bool {
		return cmpArrayOne(left, right, opFloatCmpRegex)
	})
}
