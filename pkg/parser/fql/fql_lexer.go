// Code generated from antlr/FqlLexer.g4 by ANTLR 4.9.2. DO NOT EDIT.

package fql

import (
	"fmt"
	"unicode"

	"github.com/antlr/antlr4/runtime/Go/antlr"
)

// Suppress unused import error
var _ = fmt.Printf
var _ = unicode.IsLetter

var serializedLexerAtn = []uint16{
	3, 24715, 42794, 33075, 47597, 16764, 15335, 30598, 22884, 2, 74, 623,
	8, 1, 4, 2, 9, 2, 4, 3, 9, 3, 4, 4, 9, 4, 4, 5, 9, 5, 4, 6, 9, 6, 4, 7,
	9, 7, 4, 8, 9, 8, 4, 9, 9, 9, 4, 10, 9, 10, 4, 11, 9, 11, 4, 12, 9, 12,
	4, 13, 9, 13, 4, 14, 9, 14, 4, 15, 9, 15, 4, 16, 9, 16, 4, 17, 9, 17, 4,
	18, 9, 18, 4, 19, 9, 19, 4, 20, 9, 20, 4, 21, 9, 21, 4, 22, 9, 22, 4, 23,
	9, 23, 4, 24, 9, 24, 4, 25, 9, 25, 4, 26, 9, 26, 4, 27, 9, 27, 4, 28, 9,
	28, 4, 29, 9, 29, 4, 30, 9, 30, 4, 31, 9, 31, 4, 32, 9, 32, 4, 33, 9, 33,
	4, 34, 9, 34, 4, 35, 9, 35, 4, 36, 9, 36, 4, 37, 9, 37, 4, 38, 9, 38, 4,
	39, 9, 39, 4, 40, 9, 40, 4, 41, 9, 41, 4, 42, 9, 42, 4, 43, 9, 43, 4, 44,
	9, 44, 4, 45, 9, 45, 4, 46, 9, 46, 4, 47, 9, 47, 4, 48, 9, 48, 4, 49, 9,
	49, 4, 50, 9, 50, 4, 51, 9, 51, 4, 52, 9, 52, 4, 53, 9, 53, 4, 54, 9, 54,
	4, 55, 9, 55, 4, 56, 9, 56, 4, 57, 9, 57, 4, 58, 9, 58, 4, 59, 9, 59, 4,
	60, 9, 60, 4, 61, 9, 61, 4, 62, 9, 62, 4, 63, 9, 63, 4, 64, 9, 64, 4, 65,
	9, 65, 4, 66, 9, 66, 4, 67, 9, 67, 4, 68, 9, 68, 4, 69, 9, 69, 4, 70, 9,
	70, 4, 71, 9, 71, 4, 72, 9, 72, 4, 73, 9, 73, 4, 74, 9, 74, 4, 75, 9, 75,
	4, 76, 9, 76, 4, 77, 9, 77, 4, 78, 9, 78, 4, 79, 9, 79, 4, 80, 9, 80, 4,
	81, 9, 81, 4, 82, 9, 82, 4, 83, 9, 83, 4, 84, 9, 84, 4, 85, 9, 85, 3, 2,
	3, 2, 3, 2, 3, 2, 7, 2, 176, 10, 2, 12, 2, 14, 2, 179, 11, 2, 3, 2, 3,
	2, 3, 2, 3, 2, 3, 2, 3, 3, 3, 3, 3, 3, 3, 3, 7, 3, 190, 10, 3, 12, 3, 14,
	3, 193, 11, 3, 3, 3, 3, 3, 3, 4, 6, 4, 198, 10, 4, 13, 4, 14, 4, 199, 3,
	4, 3, 4, 3, 5, 3, 5, 3, 5, 3, 5, 3, 6, 3, 6, 3, 7, 3, 7, 3, 8, 3, 8, 3,
	9, 3, 9, 3, 10, 3, 10, 3, 11, 3, 11, 3, 12, 3, 12, 3, 13, 3, 13, 3, 14,
	3, 14, 3, 15, 3, 15, 3, 16, 3, 16, 3, 17, 3, 17, 3, 18, 3, 18, 3, 18, 3,
	19, 3, 19, 3, 19, 3, 20, 3, 20, 3, 20, 3, 21, 3, 21, 3, 21, 3, 22, 3, 22,
	3, 23, 3, 23, 3, 24, 3, 24, 3, 25, 3, 25, 3, 26, 3, 26, 3, 27, 3, 27, 3,
	27, 3, 28, 3, 28, 3, 28, 3, 29, 3, 29, 3, 29, 3, 29, 3, 29, 5, 29, 265,
	10, 29, 3, 30, 3, 30, 3, 30, 3, 30, 5, 30, 271, 10, 30, 3, 31, 3, 31, 3,
	31, 3, 32, 3, 32, 3, 33, 3, 33, 3, 34, 3, 34, 3, 34, 3, 35, 3, 35, 3, 35,
	3, 36, 3, 36, 3, 36, 3, 36, 3, 37, 3, 37, 3, 37, 3, 37, 3, 37, 3, 37, 3,
	37, 3, 38, 3, 38, 3, 38, 3, 38, 3, 38, 3, 38, 3, 38, 3, 38, 3, 39, 3, 39,
	3, 39, 3, 39, 3, 39, 3, 39, 3, 39, 3, 39, 3, 40, 3, 40, 3, 40, 3, 40, 3,
	40, 3, 40, 3, 40, 3, 40, 3, 41, 3, 41, 3, 41, 3, 41, 3, 41, 3, 41, 3, 41,
	3, 41, 3, 41, 3, 42, 3, 42, 3, 42, 3, 42, 3, 42, 3, 42, 3, 42, 3, 43, 3,
	43, 3, 43, 3, 43, 3, 43, 3, 43, 3, 43, 3, 43, 3, 44, 3, 44, 3, 44, 3, 44,
	3, 44, 3, 45, 3, 45, 3, 45, 3, 45, 3, 45, 3, 45, 3, 46, 3, 46, 3, 46, 3,
	46, 3, 47, 3, 47, 3, 47, 3, 47, 3, 47, 3, 47, 3, 47, 3, 47, 3, 48, 3, 48,
	3, 48, 3, 48, 3, 48, 3, 48, 3, 48, 5, 48, 375, 10, 48, 3, 49, 3, 49, 3,
	49, 3, 49, 3, 49, 3, 50, 3, 50, 3, 50, 3, 50, 3, 50, 3, 51, 3, 51, 3, 51,
	3, 51, 3, 51, 3, 51, 3, 51, 3, 51, 3, 51, 3, 51, 3, 51, 3, 51, 3, 51, 3,
	51, 3, 51, 3, 51, 3, 51, 3, 51, 5, 51, 405, 10, 51, 3, 52, 3, 52, 3, 52,
	3, 52, 3, 53, 3, 53, 3, 53, 3, 53, 3, 53, 3, 54, 3, 54, 3, 54, 3, 54, 3,
	54, 3, 55, 3, 55, 3, 55, 3, 55, 3, 55, 3, 56, 3, 56, 3, 56, 3, 56, 3, 56,
	3, 56, 3, 57, 3, 57, 3, 57, 3, 57, 3, 58, 3, 58, 3, 58, 3, 58, 3, 59, 3,
	59, 3, 59, 3, 59, 3, 59, 3, 59, 3, 59, 3, 59, 3, 59, 3, 59, 3, 60, 3, 60,
	3, 60, 3, 60, 3, 60, 3, 60, 3, 61, 3, 61, 3, 61, 3, 61, 3, 61, 3, 62, 3,
	62, 3, 62, 3, 62, 5, 62, 465, 10, 62, 3, 63, 3, 63, 3, 63, 3, 64, 3, 64,
	3, 64, 3, 65, 3, 65, 3, 65, 3, 65, 3, 65, 3, 65, 3, 66, 3, 66, 3, 67, 6,
	67, 482, 10, 67, 13, 67, 14, 67, 483, 3, 67, 3, 67, 7, 67, 488, 10, 67,
	12, 67, 14, 67, 491, 11, 67, 7, 67, 493, 10, 67, 12, 67, 14, 67, 496, 11,
	67, 3, 67, 3, 67, 7, 67, 500, 10, 67, 12, 67, 14, 67, 503, 11, 67, 7, 67,
	505, 10, 67, 12, 67, 14, 67, 508, 11, 67, 3, 68, 3, 68, 3, 69, 3, 69, 3,
	69, 3, 69, 5, 69, 516, 10, 69, 3, 70, 6, 70, 519, 10, 70, 13, 70, 14, 70,
	520, 3, 71, 3, 71, 3, 71, 6, 71, 526, 10, 71, 13, 71, 14, 71, 527, 3, 71,
	5, 71, 531, 10, 71, 3, 71, 3, 71, 5, 71, 535, 10, 71, 5, 71, 537, 10, 71,
	3, 72, 3, 72, 3, 72, 3, 73, 3, 73, 3, 74, 3, 74, 3, 75, 3, 75, 3, 75, 7,
	75, 549, 10, 75, 12, 75, 14, 75, 552, 11, 75, 5, 75, 554, 10, 75, 3, 76,
	3, 76, 5, 76, 558, 10, 76, 3, 76, 6, 76, 561, 10, 76, 13, 76, 14, 76, 562,
	3, 77, 3, 77, 3, 78, 3, 78, 3, 79, 3, 79, 3, 80, 3, 80, 3, 81, 3, 81, 3,
	81, 3, 81, 3, 81, 3, 81, 7, 81, 579, 10, 81, 12, 81, 14, 81, 582, 11, 81,
	3, 81, 3, 81, 3, 82, 3, 82, 3, 82, 3, 82, 3, 82, 3, 82, 7, 82, 592, 10,
	82, 12, 82, 14, 82, 595, 11, 82, 3, 82, 3, 82, 3, 83, 3, 83, 3, 83, 3,
	83, 7, 83, 603, 10, 83, 12, 83, 14, 83, 606, 11, 83, 3, 83, 3, 83, 3, 84,
	3, 84, 3, 84, 3, 84, 7, 84, 614, 10, 84, 12, 84, 14, 84, 617, 11, 84, 3,
	84, 3, 84, 3, 85, 3, 85, 3, 85, 3, 177, 2, 86, 3, 3, 5, 4, 7, 5, 9, 6,
	11, 7, 13, 8, 15, 9, 17, 10, 19, 11, 21, 12, 23, 13, 25, 14, 27, 15, 29,
	16, 31, 17, 33, 18, 35, 19, 37, 20, 39, 21, 41, 22, 43, 23, 45, 24, 47,
	25, 49, 26, 51, 27, 53, 28, 55, 29, 57, 30, 59, 31, 61, 32, 63, 33, 65,
	34, 67, 35, 69, 36, 71, 37, 73, 38, 75, 39, 77, 40, 79, 41, 81, 42, 83,
	43, 85, 44, 87, 45, 89, 46, 91, 47, 93, 48, 95, 49, 97, 50, 99, 51, 101,
	52, 103, 53, 105, 54, 107, 55, 109, 56, 111, 57, 113, 58, 115, 59, 117,
	60, 119, 61, 121, 62, 123, 63, 125, 64, 127, 65, 129, 66, 131, 67, 133,
	68, 135, 69, 137, 70, 139, 71, 141, 72, 143, 73, 145, 74, 147, 2, 149,
	2, 151, 2, 153, 2, 155, 2, 157, 2, 159, 2, 161, 2, 163, 2, 165, 2, 167,
	2, 169, 2, 3, 2, 14, 5, 2, 12, 12, 15, 15, 8234, 8235, 6, 2, 11, 11, 13,
	14, 34, 34, 162, 162, 3, 2, 50, 59, 5, 2, 50, 59, 67, 72, 99, 104, 3, 2,
	51, 59, 4, 2, 71, 71, 103, 103, 4, 2, 45, 45, 47, 47, 4, 2, 67, 92, 99,
	124, 4, 2, 36, 36, 94, 94, 4, 2, 41, 41, 94, 94, 3, 2, 98, 98, 3, 2, 182,
	182, 2, 647, 2, 3, 3, 2, 2, 2, 2, 5, 3, 2, 2, 2, 2, 7, 3, 2, 2, 2, 2, 9,
	3, 2, 2, 2, 2, 11, 3, 2, 2, 2, 2, 13, 3, 2, 2, 2, 2, 15, 3, 2, 2, 2, 2,
	17, 3, 2, 2, 2, 2, 19, 3, 2, 2, 2, 2, 21, 3, 2, 2, 2, 2, 23, 3, 2, 2, 2,
	2, 25, 3, 2, 2, 2, 2, 27, 3, 2, 2, 2, 2, 29, 3, 2, 2, 2, 2, 31, 3, 2, 2,
	2, 2, 33, 3, 2, 2, 2, 2, 35, 3, 2, 2, 2, 2, 37, 3, 2, 2, 2, 2, 39, 3, 2,
	2, 2, 2, 41, 3, 2, 2, 2, 2, 43, 3, 2, 2, 2, 2, 45, 3, 2, 2, 2, 2, 47, 3,
	2, 2, 2, 2, 49, 3, 2, 2, 2, 2, 51, 3, 2, 2, 2, 2, 53, 3, 2, 2, 2, 2, 55,
	3, 2, 2, 2, 2, 57, 3, 2, 2, 2, 2, 59, 3, 2, 2, 2, 2, 61, 3, 2, 2, 2, 2,
	63, 3, 2, 2, 2, 2, 65, 3, 2, 2, 2, 2, 67, 3, 2, 2, 2, 2, 69, 3, 2, 2, 2,
	2, 71, 3, 2, 2, 2, 2, 73, 3, 2, 2, 2, 2, 75, 3, 2, 2, 2, 2, 77, 3, 2, 2,
	2, 2, 79, 3, 2, 2, 2, 2, 81, 3, 2, 2, 2, 2, 83, 3, 2, 2, 2, 2, 85, 3, 2,
	2, 2, 2, 87, 3, 2, 2, 2, 2, 89, 3, 2, 2, 2, 2, 91, 3, 2, 2, 2, 2, 93, 3,
	2, 2, 2, 2, 95, 3, 2, 2, 2, 2, 97, 3, 2, 2, 2, 2, 99, 3, 2, 2, 2, 2, 101,
	3, 2, 2, 2, 2, 103, 3, 2, 2, 2, 2, 105, 3, 2, 2, 2, 2, 107, 3, 2, 2, 2,
	2, 109, 3, 2, 2, 2, 2, 111, 3, 2, 2, 2, 2, 113, 3, 2, 2, 2, 2, 115, 3,
	2, 2, 2, 2, 117, 3, 2, 2, 2, 2, 119, 3, 2, 2, 2, 2, 121, 3, 2, 2, 2, 2,
	123, 3, 2, 2, 2, 2, 125, 3, 2, 2, 2, 2, 127, 3, 2, 2, 2, 2, 129, 3, 2,
	2, 2, 2, 131, 3, 2, 2, 2, 2, 133, 3, 2, 2, 2, 2, 135, 3, 2, 2, 2, 2, 137,
	3, 2, 2, 2, 2, 139, 3, 2, 2, 2, 2, 141, 3, 2, 2, 2, 2, 143, 3, 2, 2, 2,
	2, 145, 3, 2, 2, 2, 3, 171, 3, 2, 2, 2, 5, 185, 3, 2, 2, 2, 7, 197, 3,
	2, 2, 2, 9, 203, 3, 2, 2, 2, 11, 207, 3, 2, 2, 2, 13, 209, 3, 2, 2, 2,
	15, 211, 3, 2, 2, 2, 17, 213, 3, 2, 2, 2, 19, 215, 3, 2, 2, 2, 21, 217,
	3, 2, 2, 2, 23, 219, 3, 2, 2, 2, 25, 221, 3, 2, 2, 2, 27, 223, 3, 2, 2,
	2, 29, 225, 3, 2, 2, 2, 31, 227, 3, 2, 2, 2, 33, 229, 3, 2, 2, 2, 35, 231,
	3, 2, 2, 2, 37, 234, 3, 2, 2, 2, 39, 237, 3, 2, 2, 2, 41, 240, 3, 2, 2,
	2, 43, 243, 3, 2, 2, 2, 45, 245, 3, 2, 2, 2, 47, 247, 3, 2, 2, 2, 49, 249,
	3, 2, 2, 2, 51, 251, 3, 2, 2, 2, 53, 253, 3, 2, 2, 2, 55, 256, 3, 2, 2,
	2, 57, 264, 3, 2, 2, 2, 59, 270, 3, 2, 2, 2, 61, 272, 3, 2, 2, 2, 63, 275,
	3, 2, 2, 2, 65, 277, 3, 2, 2, 2, 67, 279, 3, 2, 2, 2, 69, 282, 3, 2, 2,
	2, 71, 285, 3, 2, 2, 2, 73, 289, 3, 2, 2, 2, 75, 296, 3, 2, 2, 2, 77, 304,
	3, 2, 2, 2, 79, 312, 3, 2, 2, 2, 81, 320, 3, 2, 2, 2, 83, 329, 3, 2, 2,
	2, 85, 336, 3, 2, 2, 2, 87, 344, 3, 2, 2, 2, 89, 349, 3, 2, 2, 2, 91, 355,
	3, 2, 2, 2, 93, 359, 3, 2, 2, 2, 95, 374, 3, 2, 2, 2, 97, 376, 3, 2, 2,
	2, 99, 381, 3, 2, 2, 2, 101, 404, 3, 2, 2, 2, 103, 406, 3, 2, 2, 2, 105,
	410, 3, 2, 2, 2, 107, 415, 3, 2, 2, 2, 109, 420, 3, 2, 2, 2, 111, 425,
	3, 2, 2, 2, 113, 431, 3, 2, 2, 2, 115, 435, 3, 2, 2, 2, 117, 439, 3, 2,
	2, 2, 119, 449, 3, 2, 2, 2, 121, 455, 3, 2, 2, 2, 123, 464, 3, 2, 2, 2,
	125, 466, 3, 2, 2, 2, 127, 469, 3, 2, 2, 2, 129, 472, 3, 2, 2, 2, 131,
	478, 3, 2, 2, 2, 133, 481, 3, 2, 2, 2, 135, 509, 3, 2, 2, 2, 137, 515,
	3, 2, 2, 2, 139, 518, 3, 2, 2, 2, 141, 536, 3, 2, 2, 2, 143, 538, 3, 2,
	2, 2, 145, 541, 3, 2, 2, 2, 147, 543, 3, 2, 2, 2, 149, 553, 3, 2, 2, 2,
	151, 555, 3, 2, 2, 2, 153, 564, 3, 2, 2, 2, 155, 566, 3, 2, 2, 2, 157,
	568, 3, 2, 2, 2, 159, 570, 3, 2, 2, 2, 161, 572, 3, 2, 2, 2, 163, 585,
	3, 2, 2, 2, 165, 598, 3, 2, 2, 2, 167, 609, 3, 2, 2, 2, 169, 620, 3, 2,
	2, 2, 171, 172, 7, 49, 2, 2, 172, 173, 7, 44, 2, 2, 173, 177, 3, 2, 2,
	2, 174, 176, 11, 2, 2, 2, 175, 174, 3, 2, 2, 2, 176, 179, 3, 2, 2, 2, 177,
	178, 3, 2, 2, 2, 177, 175, 3, 2, 2, 2, 178, 180, 3, 2, 2, 2, 179, 177,
	3, 2, 2, 2, 180, 181, 7, 44, 2, 2, 181, 182, 7, 49, 2, 2, 182, 183, 3,
	2, 2, 2, 183, 184, 8, 2, 2, 2, 184, 4, 3, 2, 2, 2, 185, 186, 7, 49, 2,
	2, 186, 187, 7, 49, 2, 2, 187, 191, 3, 2, 2, 2, 188, 190, 10, 2, 2, 2,
	189, 188, 3, 2, 2, 2, 190, 193, 3, 2, 2, 2, 191, 189, 3, 2, 2, 2, 191,
	192, 3, 2, 2, 2, 192, 194, 3, 2, 2, 2, 193, 191, 3, 2, 2, 2, 194, 195,
	8, 3, 2, 2, 195, 6, 3, 2, 2, 2, 196, 198, 9, 3, 2, 2, 197, 196, 3, 2, 2,
	2, 198, 199, 3, 2, 2, 2, 199, 197, 3, 2, 2, 2, 199, 200, 3, 2, 2, 2, 200,
	201, 3, 2, 2, 2, 201, 202, 8, 4, 2, 2, 202, 8, 3, 2, 2, 2, 203, 204, 9,
	2, 2, 2, 204, 205, 3, 2, 2, 2, 205, 206, 8, 5, 2, 2, 206, 10, 3, 2, 2,
	2, 207, 208, 7, 60, 2, 2, 208, 12, 3, 2, 2, 2, 209, 210, 7, 61, 2, 2, 210,
	14, 3, 2, 2, 2, 211, 212, 7, 48, 2, 2, 212, 16, 3, 2, 2, 2, 213, 214, 7,
	46, 2, 2, 214, 18, 3, 2, 2, 2, 215, 216, 7, 93, 2, 2, 216, 20, 3, 2, 2,
	2, 217, 218, 7, 95, 2, 2, 218, 22, 3, 2, 2, 2, 219, 220, 7, 42, 2, 2, 220,
	24, 3, 2, 2, 2, 221, 222, 7, 43, 2, 2, 222, 26, 3, 2, 2, 2, 223, 224, 7,
	125, 2, 2, 224, 28, 3, 2, 2, 2, 225, 226, 7, 127, 2, 2, 226, 30, 3, 2,
	2, 2, 227, 228, 7, 64, 2, 2, 228, 32, 3, 2, 2, 2, 229, 230, 7, 62, 2, 2,
	230, 34, 3, 2, 2, 2, 231, 232, 7, 63, 2, 2, 232, 233, 7, 63, 2, 2, 233,
	36, 3, 2, 2, 2, 234, 235, 7, 64, 2, 2, 235, 236, 7, 63, 2, 2, 236, 38,
	3, 2, 2, 2, 237, 238, 7, 62, 2, 2, 238, 239, 7, 63, 2, 2, 239, 40, 3, 2,
	2, 2, 240, 241, 7, 35, 2, 2, 241, 242, 7, 63, 2, 2, 242, 42, 3, 2, 2, 2,
	243, 244, 7, 44, 2, 2, 244, 44, 3, 2, 2, 2, 245, 246, 7, 49, 2, 2, 246,
	46, 3, 2, 2, 2, 247, 248, 7, 39, 2, 2, 248, 48, 3, 2, 2, 2, 249, 250, 7,
	45, 2, 2, 250, 50, 3, 2, 2, 2, 251, 252, 7, 47, 2, 2, 252, 52, 3, 2, 2,
	2, 253, 254, 7, 47, 2, 2, 254, 255, 7, 47, 2, 2, 255, 54, 3, 2, 2, 2, 256,
	257, 7, 45, 2, 2, 257, 258, 7, 45, 2, 2, 258, 56, 3, 2, 2, 2, 259, 260,
	7, 67, 2, 2, 260, 261, 7, 80, 2, 2, 261, 265, 7, 70, 2, 2, 262, 263, 7,
	40, 2, 2, 263, 265, 7, 40, 2, 2, 264, 259, 3, 2, 2, 2, 264, 262, 3, 2,
	2, 2, 265, 58, 3, 2, 2, 2, 266, 267, 7, 81, 2, 2, 267, 271, 7, 84, 2, 2,
	268, 269, 7, 126, 2, 2, 269, 271, 7, 126, 2, 2, 270, 266, 3, 2, 2, 2, 270,
	268, 3, 2, 2, 2, 271, 60, 3, 2, 2, 2, 272, 273, 5, 15, 8, 2, 273, 274,
	5, 15, 8, 2, 274, 62, 3, 2, 2, 2, 275, 276, 7, 63, 2, 2, 276, 64, 3, 2,
	2, 2, 277, 278, 7, 65, 2, 2, 278, 66, 3, 2, 2, 2, 279, 280, 7, 35, 2, 2,
	280, 281, 7, 128, 2, 2, 281, 68, 3, 2, 2, 2, 282, 283, 7, 63, 2, 2, 283,
	284, 7, 128, 2, 2, 284, 70, 3, 2, 2, 2, 285, 286, 7, 72, 2, 2, 286, 287,
	7, 81, 2, 2, 287, 288, 7, 84, 2, 2, 288, 72, 3, 2, 2, 2, 289, 290, 7, 84,
	2, 2, 290, 291, 7, 71, 2, 2, 291, 292, 7, 86, 2, 2, 292, 293, 7, 87, 2,
	2, 293, 294, 7, 84, 2, 2, 294, 295, 7, 80, 2, 2, 295, 74, 3, 2, 2, 2, 296,
	297, 7, 89, 2, 2, 297, 298, 7, 67, 2, 2, 298, 299, 7, 75, 2, 2, 299, 300,
	7, 86, 2, 2, 300, 301, 7, 72, 2, 2, 301, 302, 7, 81, 2, 2, 302, 303, 7,
	84, 2, 2, 303, 76, 3, 2, 2, 2, 304, 305, 7, 81, 2, 2, 305, 306, 7, 82,
	2, 2, 306, 307, 7, 86, 2, 2, 307, 308, 7, 75, 2, 2, 308, 309, 7, 81, 2,
	2, 309, 310, 7, 80, 2, 2, 310, 311, 7, 85, 2, 2, 311, 78, 3, 2, 2, 2, 312,
	313, 7, 86, 2, 2, 313, 314, 7, 75, 2, 2, 314, 315, 7, 79, 2, 2, 315, 316,
	7, 71, 2, 2, 316, 317, 7, 81, 2, 2, 317, 318, 7, 87, 2, 2, 318, 319, 7,
	86, 2, 2, 319, 80, 3, 2, 2, 2, 320, 321, 7, 70, 2, 2, 321, 322, 7, 75,
	2, 2, 322, 323, 7, 85, 2, 2, 323, 324, 7, 86, 2, 2, 324, 325, 7, 75, 2,
	2, 325, 326, 7, 80, 2, 2, 326, 327, 7, 69, 2, 2, 327, 328, 7, 86, 2, 2,
	328, 82, 3, 2, 2, 2, 329, 330, 7, 72, 2, 2, 330, 331, 7, 75, 2, 2, 331,
	332, 7, 78, 2, 2, 332, 333, 7, 86, 2, 2, 333, 334, 7, 71, 2, 2, 334, 335,
	7, 84, 2, 2, 335, 84, 3, 2, 2, 2, 336, 337, 7, 69, 2, 2, 337, 338, 7, 87,
	2, 2, 338, 339, 7, 84, 2, 2, 339, 340, 7, 84, 2, 2, 340, 341, 7, 71, 2,
	2, 341, 342, 7, 80, 2, 2, 342, 343, 7, 86, 2, 2, 343, 86, 3, 2, 2, 2, 344,
	345, 7, 85, 2, 2, 345, 346, 7, 81, 2, 2, 346, 347, 7, 84, 2, 2, 347, 348,
	7, 86, 2, 2, 348, 88, 3, 2, 2, 2, 349, 350, 7, 78, 2, 2, 350, 351, 7, 75,
	2, 2, 351, 352, 7, 79, 2, 2, 352, 353, 7, 75, 2, 2, 353, 354, 7, 86, 2,
	2, 354, 90, 3, 2, 2, 2, 355, 356, 7, 78, 2, 2, 356, 357, 7, 71, 2, 2, 357,
	358, 7, 86, 2, 2, 358, 92, 3, 2, 2, 2, 359, 360, 7, 69, 2, 2, 360, 361,
	7, 81, 2, 2, 361, 362, 7, 78, 2, 2, 362, 363, 7, 78, 2, 2, 363, 364, 7,
	71, 2, 2, 364, 365, 7, 69, 2, 2, 365, 366, 7, 86, 2, 2, 366, 94, 3, 2,
	2, 2, 367, 368, 7, 67, 2, 2, 368, 369, 7, 85, 2, 2, 369, 375, 7, 69, 2,
	2, 370, 371, 7, 70, 2, 2, 371, 372, 7, 71, 2, 2, 372, 373, 7, 85, 2, 2,
	373, 375, 7, 69, 2, 2, 374, 367, 3, 2, 2, 2, 374, 370, 3, 2, 2, 2, 375,
	96, 3, 2, 2, 2, 376, 377, 7, 80, 2, 2, 377, 378, 7, 81, 2, 2, 378, 379,
	7, 80, 2, 2, 379, 380, 7, 71, 2, 2, 380, 98, 3, 2, 2, 2, 381, 382, 7, 80,
	2, 2, 382, 383, 7, 87, 2, 2, 383, 384, 7, 78, 2, 2, 384, 385, 7, 78, 2,
	2, 385, 100, 3, 2, 2, 2, 386, 387, 7, 86, 2, 2, 387, 388, 7, 84, 2, 2,
	388, 389, 7, 87, 2, 2, 389, 405, 7, 71, 2, 2, 390, 391, 7, 118, 2, 2, 391,
	392, 7, 116, 2, 2, 392, 393, 7, 119, 2, 2, 393, 405, 7, 103, 2, 2, 394,
	395, 7, 72, 2, 2, 395, 396, 7, 67, 2, 2, 396, 397, 7, 78, 2, 2, 397, 398,
	7, 85, 2, 2, 398, 405, 7, 71, 2, 2, 399, 400, 7, 104, 2, 2, 400, 401, 7,
	99, 2, 2, 401, 402, 7, 110, 2, 2, 402, 403, 7, 117, 2, 2, 403, 405, 7,
	103, 2, 2, 404, 386, 3, 2, 2, 2, 404, 390, 3, 2, 2, 2, 404, 394, 3, 2,
	2, 2, 404, 399, 3, 2, 2, 2, 405, 102, 3, 2, 2, 2, 406, 407, 7, 87, 2, 2,
	407, 408, 7, 85, 2, 2, 408, 409, 7, 71, 2, 2, 409, 104, 3, 2, 2, 2, 410,
	411, 7, 75, 2, 2, 411, 412, 7, 80, 2, 2, 412, 413, 7, 86, 2, 2, 413, 414,
	7, 81, 2, 2, 414, 106, 3, 2, 2, 2, 415, 416, 7, 77, 2, 2, 416, 417, 7,
	71, 2, 2, 417, 418, 7, 71, 2, 2, 418, 419, 7, 82, 2, 2, 419, 108, 3, 2,
	2, 2, 420, 421, 7, 89, 2, 2, 421, 422, 7, 75, 2, 2, 422, 423, 7, 86, 2,
	2, 423, 424, 7, 74, 2, 2, 424, 110, 3, 2, 2, 2, 425, 426, 7, 69, 2, 2,
	426, 427, 7, 81, 2, 2, 427, 428, 7, 87, 2, 2, 428, 429, 7, 80, 2, 2, 429,
	430, 7, 86, 2, 2, 430, 112, 3, 2, 2, 2, 431, 432, 7, 67, 2, 2, 432, 433,
	7, 78, 2, 2, 433, 434, 7, 78, 2, 2, 434, 114, 3, 2, 2, 2, 435, 436, 7,
	67, 2, 2, 436, 437, 7, 80, 2, 2, 437, 438, 7, 91, 2, 2, 438, 116, 3, 2,
	2, 2, 439, 440, 7, 67, 2, 2, 440, 441, 7, 73, 2, 2, 441, 442, 7, 73, 2,
	2, 442, 443, 7, 84, 2, 2, 443, 444, 7, 71, 2, 2, 444, 445, 7, 73, 2, 2,
	445, 446, 7, 67, 2, 2, 446, 447, 7, 86, 2, 2, 447, 448, 7, 71, 2, 2, 448,
	118, 3, 2, 2, 2, 449, 450, 7, 71, 2, 2, 450, 451, 7, 88, 2, 2, 451, 452,
	7, 71, 2, 2, 452, 453, 7, 80, 2, 2, 453, 454, 7, 86, 2, 2, 454, 120, 3,
	2, 2, 2, 455, 456, 7, 78, 2, 2, 456, 457, 7, 75, 2, 2, 457, 458, 7, 77,
	2, 2, 458, 459, 7, 71, 2, 2, 459, 122, 3, 2, 2, 2, 460, 461, 7, 80, 2,
	2, 461, 462, 7, 81, 2, 2, 462, 465, 7, 86, 2, 2, 463, 465, 7, 35, 2, 2,
	464, 460, 3, 2, 2, 2, 464, 463, 3, 2, 2, 2, 465, 124, 3, 2, 2, 2, 466,
	467, 7, 75, 2, 2, 467, 468, 7, 80, 2, 2, 468, 126, 3, 2, 2, 2, 469, 470,
	7, 70, 2, 2, 470, 471, 7, 81, 2, 2, 471, 128, 3, 2, 2, 2, 472, 473, 7,
	89, 2, 2, 473, 474, 7, 74, 2, 2, 474, 475, 7, 75, 2, 2, 475, 476, 7, 78,
	2, 2, 476, 477, 7, 71, 2, 2, 477, 130, 3, 2, 2, 2, 478, 479, 7, 66, 2,
	2, 479, 132, 3, 2, 2, 2, 480, 482, 5, 153, 77, 2, 481, 480, 3, 2, 2, 2,
	482, 483, 3, 2, 2, 2, 483, 481, 3, 2, 2, 2, 483, 484, 3, 2, 2, 2, 484,
	494, 3, 2, 2, 2, 485, 489, 5, 155, 78, 2, 486, 488, 5, 133, 67, 2, 487,
	486, 3, 2, 2, 2, 488, 491, 3, 2, 2, 2, 489, 487, 3, 2, 2, 2, 489, 490,
	3, 2, 2, 2, 490, 493, 3, 2, 2, 2, 491, 489, 3, 2, 2, 2, 492, 485, 3, 2,
	2, 2, 493, 496, 3, 2, 2, 2, 494, 492, 3, 2, 2, 2, 494, 495, 3, 2, 2, 2,
	495, 506, 3, 2, 2, 2, 496, 494, 3, 2, 2, 2, 497, 501, 5, 159, 80, 2, 498,
	500, 5, 133, 67, 2, 499, 498, 3, 2, 2, 2, 500, 503, 3, 2, 2, 2, 501, 499,
	3, 2, 2, 2, 501, 502, 3, 2, 2, 2, 502, 505, 3, 2, 2, 2, 503, 501, 3, 2,
	2, 2, 504, 497, 3, 2, 2, 2, 505, 508, 3, 2, 2, 2, 506, 504, 3, 2, 2, 2,
	506, 507, 3, 2, 2, 2, 507, 134, 3, 2, 2, 2, 508, 506, 3, 2, 2, 2, 509,
	510, 5, 157, 79, 2, 510, 136, 3, 2, 2, 2, 511, 516, 5, 163, 82, 2, 512,
	516, 5, 161, 81, 2, 513, 516, 5, 165, 83, 2, 514, 516, 5, 167, 84, 2, 515,
	511, 3, 2, 2, 2, 515, 512, 3, 2, 2, 2, 515, 513, 3, 2, 2, 2, 515, 514,
	3, 2, 2, 2, 516, 138, 3, 2, 2, 2, 517, 519, 9, 4, 2, 2, 518, 517, 3, 2,
	2, 2, 519, 520, 3, 2, 2, 2, 520, 518, 3, 2, 2, 2, 520, 521, 3, 2, 2, 2,
	521, 140, 3, 2, 2, 2, 522, 523, 5, 149, 75, 2, 523, 525, 5, 15, 8, 2, 524,
	526, 9, 4, 2, 2, 525, 524, 3, 2, 2, 2, 526, 527, 3, 2, 2, 2, 527, 525,
	3, 2, 2, 2, 527, 528, 3, 2, 2, 2, 528, 530, 3, 2, 2, 2, 529, 531, 5, 151,
	76, 2, 530, 529, 3, 2, 2, 2, 530, 531, 3, 2, 2, 2, 531, 537, 3, 2, 2, 2,
	532, 534, 5, 149, 75, 2, 533, 535, 5, 151, 76, 2, 534, 533, 3, 2, 2, 2,
	534, 535, 3, 2, 2, 2, 535, 537, 3, 2, 2, 2, 536, 522, 3, 2, 2, 2, 536,
	532, 3, 2, 2, 2, 537, 142, 3, 2, 2, 2, 538, 539, 5, 133, 67, 2, 539, 540,
	5, 169, 85, 2, 540, 144, 3, 2, 2, 2, 541, 542, 11, 2, 2, 2, 542, 146, 3,
	2, 2, 2, 543, 544, 9, 5, 2, 2, 544, 148, 3, 2, 2, 2, 545, 554, 7, 50, 2,
	2, 546, 550, 9, 6, 2, 2, 547, 549, 9, 4, 2, 2, 548, 547, 3, 2, 2, 2, 549,
	552, 3, 2, 2, 2, 550, 548, 3, 2, 2, 2, 550, 551, 3, 2, 2, 2, 551, 554,
	3, 2, 2, 2, 552, 550, 3, 2, 2, 2, 553, 545, 3, 2, 2, 2, 553, 546, 3, 2,
	2, 2, 554, 150, 3, 2, 2, 2, 555, 557, 9, 7, 2, 2, 556, 558, 9, 8, 2, 2,
	557, 556, 3, 2, 2, 2, 557, 558, 3, 2, 2, 2, 558, 560, 3, 2, 2, 2, 559,
	561, 9, 4, 2, 2, 560, 559, 3, 2, 2, 2, 561, 562, 3, 2, 2, 2, 562, 560,
	3, 2, 2, 2, 562, 563, 3, 2, 2, 2, 563, 152, 3, 2, 2, 2, 564, 565, 9, 9,
	2, 2, 565, 154, 3, 2, 2, 2, 566, 567, 5, 157, 79, 2, 567, 156, 3, 2, 2,
	2, 568, 569, 7, 97, 2, 2, 569, 158, 3, 2, 2, 2, 570, 571, 4, 50, 59, 2,
	571, 160, 3, 2, 2, 2, 572, 580, 7, 36, 2, 2, 573, 574, 7, 94, 2, 2, 574,
	579, 11, 2, 2, 2, 575, 576, 7, 36, 2, 2, 576, 579, 7, 36, 2, 2, 577, 579,
	10, 10, 2, 2, 578, 573, 3, 2, 2, 2, 578, 575, 3, 2, 2, 2, 578, 577, 3,
	2, 2, 2, 579, 582, 3, 2, 2, 2, 580, 578, 3, 2, 2, 2, 580, 581, 3, 2, 2,
	2, 581, 583, 3, 2, 2, 2, 582, 580, 3, 2, 2, 2, 583, 584, 7, 36, 2, 2, 584,
	162, 3, 2, 2, 2, 585, 593, 7, 41, 2, 2, 586, 587, 7, 94, 2, 2, 587, 592,
	11, 2, 2, 2, 588, 589, 7, 41, 2, 2, 589, 592, 7, 41, 2, 2, 590, 592, 10,
	11, 2, 2, 591, 586, 3, 2, 2, 2, 591, 588, 3, 2, 2, 2, 591, 590, 3, 2, 2,
	2, 592, 595, 3, 2, 2, 2, 593, 591, 3, 2, 2, 2, 593, 594, 3, 2, 2, 2, 594,
	596, 3, 2, 2, 2, 595, 593, 3, 2, 2, 2, 596, 597, 7, 41, 2, 2, 597, 164,
	3, 2, 2, 2, 598, 604, 7, 98, 2, 2, 599, 600, 7, 94, 2, 2, 600, 603, 7,
	98, 2, 2, 601, 603, 10, 12, 2, 2, 602, 599, 3, 2, 2, 2, 602, 601, 3, 2,
	2, 2, 603, 606, 3, 2, 2, 2, 604, 602, 3, 2, 2, 2, 604, 605, 3, 2, 2, 2,
	605, 607, 3, 2, 2, 2, 606, 604, 3, 2, 2, 2, 607, 608, 7, 98, 2, 2, 608,
	166, 3, 2, 2, 2, 609, 615, 7, 182, 2, 2, 610, 611, 7, 94, 2, 2, 611, 614,
	7, 182, 2, 2, 612, 614, 10, 13, 2, 2, 613, 610, 3, 2, 2, 2, 613, 612, 3,
	2, 2, 2, 614, 617, 3, 2, 2, 2, 615, 613, 3, 2, 2, 2, 615, 616, 3, 2, 2,
	2, 616, 618, 3, 2, 2, 2, 617, 615, 3, 2, 2, 2, 618, 619, 7, 182, 2, 2,
	619, 168, 3, 2, 2, 2, 620, 621, 7, 60, 2, 2, 621, 622, 7, 60, 2, 2, 622,
	170, 3, 2, 2, 2, 34, 2, 177, 191, 199, 264, 270, 374, 404, 464, 483, 489,
	494, 501, 506, 515, 520, 527, 530, 534, 536, 550, 553, 557, 562, 578, 580,
	591, 593, 602, 604, 613, 615, 3, 2, 3, 2,
}

var lexerChannelNames = []string{
	"DEFAULT_TOKEN_CHANNEL", "HIDDEN",
}

var lexerModeNames = []string{
	"DEFAULT_MODE",
}

var lexerLiteralNames = []string{
	"", "", "", "", "", "':'", "';'", "'.'", "','", "'['", "']'", "'('", "')'",
	"'{'", "'}'", "'>'", "'<'", "'=='", "'>='", "'<='", "'!='", "'*'", "'/'",
	"'%'", "'+'", "'-'", "'--'", "'++'", "", "", "", "'='", "'?'", "'!~'",
	"'=~'", "'FOR'", "'RETURN'", "'WAITFOR'", "'OPTIONS'", "'TIMEOUT'", "'DISTINCT'",
	"'FILTER'", "'CURRENT'", "'SORT'", "'LIMIT'", "'LET'", "'COLLECT'", "",
	"'NONE'", "'NULL'", "", "'USE'", "'INTO'", "'KEEP'", "'WITH'", "'COUNT'",
	"'ALL'", "'ANY'", "'AGGREGATE'", "'EVENT'", "'LIKE'", "", "'IN'", "'DO'",
	"'WHILE'", "'@'",
}

var lexerSymbolicNames = []string{
	"", "MultiLineComment", "SingleLineComment", "WhiteSpaces", "LineTerminator",
	"Colon", "SemiColon", "Dot", "Comma", "OpenBracket", "CloseBracket", "OpenParen",
	"CloseParen", "OpenBrace", "CloseBrace", "Gt", "Lt", "Eq", "Gte", "Lte",
	"Neq", "Multi", "Div", "Mod", "Plus", "Minus", "MinusMinus", "PlusPlus",
	"And", "Or", "Range", "Assign", "QuestionMark", "RegexNotMatch", "RegexMatch",
	"For", "Return", "Waitfor", "Options", "Timeout", "Distinct", "Filter",
	"Current", "Sort", "Limit", "Let", "Collect", "SortDirection", "None",
	"Null", "BooleanLiteral", "Use", "Into", "Keep", "With", "Count", "All",
	"Any", "Aggregate", "Event", "Like", "Not", "In", "Do", "While", "Param",
	"Identifier", "IgnoreIdentifier", "StringLiteral", "IntegerLiteral", "FloatLiteral",
	"NamespaceSegment", "UnknownIdentifier",
}

var lexerRuleNames = []string{
	"MultiLineComment", "SingleLineComment", "WhiteSpaces", "LineTerminator",
	"Colon", "SemiColon", "Dot", "Comma", "OpenBracket", "CloseBracket", "OpenParen",
	"CloseParen", "OpenBrace", "CloseBrace", "Gt", "Lt", "Eq", "Gte", "Lte",
	"Neq", "Multi", "Div", "Mod", "Plus", "Minus", "MinusMinus", "PlusPlus",
	"And", "Or", "Range", "Assign", "QuestionMark", "RegexNotMatch", "RegexMatch",
	"For", "Return", "Waitfor", "Options", "Timeout", "Distinct", "Filter",
	"Current", "Sort", "Limit", "Let", "Collect", "SortDirection", "None",
	"Null", "BooleanLiteral", "Use", "Into", "Keep", "With", "Count", "All",
	"Any", "Aggregate", "Event", "Like", "Not", "In", "Do", "While", "Param",
	"Identifier", "IgnoreIdentifier", "StringLiteral", "IntegerLiteral", "FloatLiteral",
	"NamespaceSegment", "UnknownIdentifier", "HexDigit", "DecimalIntegerLiteral",
	"ExponentPart", "Letter", "Symbols", "Underscore", "Digit", "DQSring",
	"SQString", "BacktickString", "TickString", "NamespaceSeparator",
}

type FqlLexer struct {
	*antlr.BaseLexer
	channelNames []string
	modeNames    []string
	// TODO: EOF string
}

// NewFqlLexer produces a new lexer instance for the optional input antlr.CharStream.
//
// The *FqlLexer instance produced may be reused by calling the SetInputStream method.
// The initial lexer configuration is expensive to construct, and the object is not thread-safe;
// however, if used within a Golang sync.Pool, the construction cost amortizes well and the
// objects can be used in a thread-safe manner.
func NewFqlLexer(input antlr.CharStream) *FqlLexer {
	l := new(FqlLexer)
	lexerDeserializer := antlr.NewATNDeserializer(nil)
	lexerAtn := lexerDeserializer.DeserializeFromUInt16(serializedLexerAtn)
	lexerDecisionToDFA := make([]*antlr.DFA, len(lexerAtn.DecisionToState))
	for index, ds := range lexerAtn.DecisionToState {
		lexerDecisionToDFA[index] = antlr.NewDFA(ds, index)
	}
	l.BaseLexer = antlr.NewBaseLexer(input)
	l.Interpreter = antlr.NewLexerATNSimulator(l, lexerAtn, lexerDecisionToDFA, antlr.NewPredictionContextCache())

	l.channelNames = lexerChannelNames
	l.modeNames = lexerModeNames
	l.RuleNames = lexerRuleNames
	l.LiteralNames = lexerLiteralNames
	l.SymbolicNames = lexerSymbolicNames
	l.GrammarFileName = "FqlLexer.g4"
	// TODO: l.EOF = antlr.TokenEOF

	return l
}

// FqlLexer tokens.
const (
	FqlLexerMultiLineComment  = 1
	FqlLexerSingleLineComment = 2
	FqlLexerWhiteSpaces       = 3
	FqlLexerLineTerminator    = 4
	FqlLexerColon             = 5
	FqlLexerSemiColon         = 6
	FqlLexerDot               = 7
	FqlLexerComma             = 8
	FqlLexerOpenBracket       = 9
	FqlLexerCloseBracket      = 10
	FqlLexerOpenParen         = 11
	FqlLexerCloseParen        = 12
	FqlLexerOpenBrace         = 13
	FqlLexerCloseBrace        = 14
	FqlLexerGt                = 15
	FqlLexerLt                = 16
	FqlLexerEq                = 17
	FqlLexerGte               = 18
	FqlLexerLte               = 19
	FqlLexerNeq               = 20
	FqlLexerMulti             = 21
	FqlLexerDiv               = 22
	FqlLexerMod               = 23
	FqlLexerPlus              = 24
	FqlLexerMinus             = 25
	FqlLexerMinusMinus        = 26
	FqlLexerPlusPlus          = 27
	FqlLexerAnd               = 28
	FqlLexerOr                = 29
	FqlLexerRange             = 30
	FqlLexerAssign            = 31
	FqlLexerQuestionMark      = 32
	FqlLexerRegexNotMatch     = 33
	FqlLexerRegexMatch        = 34
	FqlLexerFor               = 35
	FqlLexerReturn            = 36
	FqlLexerWaitfor           = 37
	FqlLexerOptions           = 38
	FqlLexerTimeout           = 39
	FqlLexerDistinct          = 40
	FqlLexerFilter            = 41
	FqlLexerCurrent           = 42
	FqlLexerSort              = 43
	FqlLexerLimit             = 44
	FqlLexerLet               = 45
	FqlLexerCollect           = 46
	FqlLexerSortDirection     = 47
	FqlLexerNone              = 48
	FqlLexerNull              = 49
	FqlLexerBooleanLiteral    = 50
	FqlLexerUse               = 51
	FqlLexerInto              = 52
	FqlLexerKeep              = 53
	FqlLexerWith              = 54
	FqlLexerCount             = 55
	FqlLexerAll               = 56
	FqlLexerAny               = 57
	FqlLexerAggregate         = 58
	FqlLexerEvent             = 59
	FqlLexerLike              = 60
	FqlLexerNot               = 61
	FqlLexerIn                = 62
	FqlLexerDo                = 63
	FqlLexerWhile             = 64
	FqlLexerParam             = 65
	FqlLexerIdentifier        = 66
	FqlLexerIgnoreIdentifier  = 67
	FqlLexerStringLiteral     = 68
	FqlLexerIntegerLiteral    = 69
	FqlLexerFloatLiteral      = 70
	FqlLexerNamespaceSegment  = 71
	FqlLexerUnknownIdentifier = 72
)
