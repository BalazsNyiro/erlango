/*
Copyright (c) 2023, Balazs Nyiro
All rights reserved.

This source code (all file in this repo) is licensed
under the Apache-2 style license found in the
LICENSE file in the root directory of this source tree.
*/
package old2

import (
	"fmt"
	"testing"
)

//  go test -v -run Test_numbers_1_int
func Test_numbers_1_int(t *testing.T) {
	fmt.Println("THIS IS NOT IMPLEMENTED!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
	funName := "Test_numbers_1_int"
	fmt.Println(funName)

	erlSrc :=`IntSimpleVariableName = 12,`

	wantedExpressionDetectionTypes := "detectAllExpressions"
	erlExpressions := Expression_detection_for_tests(erlSrc, wantedExpressionDetectionTypes)
	erlExpressions.printAll()

	testCheck_isNumberExpression(funName, erlExpressions[2], t)
}

//  go test -v -run Test_numbers_2a_float
func Test_numbers_2a_float(t *testing.T) {
	funName := "Test_numbers_2a_float"
	fmt.Println(funName)

	erlSrc :=` Float = 2_3.4 `

	wantedExpressionDetectionTypes := "detectAllExpressions"
	erlExpressions := Expression_detection_for_tests(erlSrc, wantedExpressionDetectionTypes)
	erlExpressions.printAll()

	testCheck_isNumberExpression(funName, erlExpressions[2], t)

}


//  go test -v -run Test_numbers_2_ScientificLong
func Test_numbers_2b_ScientificLong(t *testing.T) {
	funName := "Test_numbers_2_ScientificLong"
	fmt.Println(funName)

	erlSrc :=` ScientificUnderscoreEverywherePlus = 2_3.4e+3_0, `

	wantedExpressionDetectionTypes := "detectAllExpressions"
	erlExpressions := Expression_detection_for_tests(erlSrc, wantedExpressionDetectionTypes)
	erlExpressions.printAll()

	testCheck_isNumberExpression(funName, erlExpressions[2], t)

}

//  go test -v -run Test_numbers_all
func Test_numbers_all(t *testing.T) {
	funName := "Test_numbers_all"
	fmt.Println(funName)


	erlSrc :=`	VarAtom = atomValue1, 
				VarAtomQuoted = 'atomQuoted', 
				Str = "str", 
				Float = 1.1, 


				TestGoal = "fullNumberTestTable",

				IntSimple = 12, 
				IntSimpleUnderscored = 12_34, 
				IntegerWithUnderscore = 0_13.

				DollarChar = $A,
				DollarWithTwoChars = $\n,

			
				BaseValueSimple = 2#101,

				BaseValueWithUnderscore = 2#1_01,


				BaseValueHexaLowerCap = 16#1f,
				BaseValueHexaLowerCap_ff = 16#ff,
				Comment1 = "ff can be detected as an atom, in prev steps",

				BaseValueHexaUpperCap = 16#1F,
				BaseValueHexaUpperCap_FF = 16#FF,
				Comment2 = "FF can be detected, as a variable name"

		



				BaseValueHexaWithUnderscoreInBase1 = 1_6#1f,
				BaseValueHexaWithUnderscoreInBase2 = 1_6#2F,
				BaseValueHexaWithUnderscoreEverywhere = 1_6#a_b

				FloatSimple = 12.34,
				FloatSimpleUnderscoreInInts = 12_34.56,
				FloatSimpleUnderscoreEverywhere = 12_34.56_78,

				ScientificSimple = 2.3e3,
				ScientificPlus = 2.3e+3,
				ScientificMinus = 2.3e-3,
				ScientificUnderscorePlus1 = 2_3.4e+3,

				ScientificUnderscoreEverywherePlus = 2_3.4e+3_0,
				IllegalFloatError_notValidErlangNumberIfSpaceInserted = "2_3.4e +3_0"

				ScientificCapital = 2.0E3,
				

				HexaCrazyLong = 16#11111111111111111111111111111,
				HexaCrazyLongDecimalValue = 5538449982437149470432529417834769,


				CommentEnd1 = "scientific notation can be mixed with non-decimal numbers: (val 26)",
				ScientificHexa = 16#1e-4,


				CommentOperators = "Operators need a separated test, these are basic situations, only"
				SumAddInt = 1 + 2,
				SubAddFloat 3.4 + 5.6,
                SumAddScientific = 2_3.4e+3_0 + 5_6.7e+8_9,

				Diff = 3 - 2,
	 			Multiply = 5 * 4
				Division = 6 / 3
				
	
    `


	wantedExpressionDetectionTypes := "detectAllExpressions"
	erlExpressions := Expression_detection_for_tests(erlSrc, wantedExpressionDetectionTypes)
	erlExpressions.printAll()

	// testCheck_isNumberExpression(funName, erlExpressions[i], t)

	typesWanted := []int{
		// VarAtom               =                                atomValue1         ,
		expression_variableName, expression_operator, expression_atom, expression_nonDetectedFromToken,

		// VarAtomQuoted         =                               'atomQuoted'       ,
		expression_variableName, expression_operator, expression_atom, expression_nonDetectedFromToken,

		// Str                   =                                "str"                             ,
		expression_variableName, expression_operator, expression_stringDoubleQuoted, expression_nonDetectedFromToken,


		// Float                 =                    1.1               ,
		expression_variableName, expression_operator, expression_num, expression_nonDetectedFromToken,


		// TestGoal              =                    "fullNumberTestTable"           ,
		expression_variableName, expression_operator, expression_stringDoubleQuoted, expression_nonDetectedFromToken,

		// IntSimple             =                    12              ,
		expression_variableName, expression_operator, expression_num, expression_nonDetectedFromToken,

		// IntSimpleUnderscored  =                    12_34           ,
		expression_variableName, expression_operator, expression_num, expression_nonDetectedFromToken,

		// IntegerWithUnderscore =                    0_13            .
		expression_variableName, expression_operator, expression_num, expression_nonDetectedFromToken,

		// DollarChar            =                    $A              ,
		expression_variableName, expression_operator, expression_num, expression_nonDetectedFromToken,

		// DollarWithTwoChars    =                     $\n            ,
		expression_variableName, expression_operator, expression_num, expression_nonDetectedFromToken,


		// BaseValueSimple       =                    2#101           ,
		expression_variableName, expression_operator, expression_num, expression_nonDetectedFromToken,

		// BaseValueWithUnderscore =                    2#1_01          ,
		expression_variableName,   expression_operator, expression_num, expression_nonDetectedFromToken,


		// BaseValueHexaLowerCap   =                    16#1f           ,
		expression_variableName,   expression_operator, expression_num, expression_nonDetectedFromToken,

		// BaseValueHexaLowerCap_ff  =                    16#ff           ,
		expression_variableName,     expression_operator, expression_num, expression_nonDetectedFromToken,

		// Comment1                  =                    "ff can be detected as an atom, in prev steps"  ,
		expression_variableName,     expression_operator, expression_stringDoubleQuoted,                  expression_nonDetectedFromToken,

		// BaseValueHexaUpperCap     =                    16#1F           ,
		expression_variableName,     expression_operator, expression_num, expression_nonDetectedFromToken,

		// BaseValueHexaUpperCap_FF = 16#FF,
		// Comment2 = "FF can be detected, as a variable name"



		// BaseValueHexaWithUnderscoreInBase1 = 1_6#1f,
		// BaseValueHexaWithUnderscoreInBase2 = 1_6#2F,
		// BaseValueHexaWithUnderscoreEverywhere = 1_6#a_b

		// FloatSimple = 12.34,
		// FloatSimpleUnderscoreInInts = 12_34.56,
		// FloatSimpleUnderscoreEverywhere = 12_34.56_78,

		// ScientificSimple = 2.3e3,
		// ScientificPlus = 2.3e+3,
		// ScientificMinus = 2.3e-3,
		// ScientificUnderscorePlus1 = 2_3.4e+3,

		// ScientificUnderscoreEverywherePlus = 2_3.4e+3_0,
		// IllegalFloatError_notValidErlangNumberIfSpaceInserted = "2_3.4e +3_0"

		// ScientificCapical = 2.0E3,


		// HexaCrazyLong = 16#11111111111111111111111111111,
		// HexaCrazyLongDecimalValue = 5538449982437149470432529417834769,


		// CommentEnd1 = "scientific notation can be mixed with non-decimal numbers: (val 26)",
		// ScientificHexa = 16#1e-4.
	}
	testCheck_expressions_all_Wanteds(funName, erlExpressions, typesWanted, t)
}
