// Code generated from antlr/FqlParser.g4 by ANTLR 4.9. DO NOT EDIT.

package fql // FqlParser
import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/antlr/antlr4/runtime/Go/antlr"
)

// Suppress unused import errors
var _ = fmt.Printf
var _ = reflect.Copy
var _ = strconv.Itoa

var parserATN = []uint16{
	3, 24715, 42794, 33075, 47597, 16764, 15335, 30598, 22884, 3, 69, 666,
	4, 2, 9, 2, 4, 3, 9, 3, 4, 4, 9, 4, 4, 5, 9, 5, 4, 6, 9, 6, 4, 7, 9, 7,
	4, 8, 9, 8, 4, 9, 9, 9, 4, 10, 9, 10, 4, 11, 9, 11, 4, 12, 9, 12, 4, 13,
	9, 13, 4, 14, 9, 14, 4, 15, 9, 15, 4, 16, 9, 16, 4, 17, 9, 17, 4, 18, 9,
	18, 4, 19, 9, 19, 4, 20, 9, 20, 4, 21, 9, 21, 4, 22, 9, 22, 4, 23, 9, 23,
	4, 24, 9, 24, 4, 25, 9, 25, 4, 26, 9, 26, 4, 27, 9, 27, 4, 28, 9, 28, 4,
	29, 9, 29, 4, 30, 9, 30, 4, 31, 9, 31, 4, 32, 9, 32, 4, 33, 9, 33, 4, 34,
	9, 34, 4, 35, 9, 35, 4, 36, 9, 36, 4, 37, 9, 37, 4, 38, 9, 38, 4, 39, 9,
	39, 4, 40, 9, 40, 4, 41, 9, 41, 4, 42, 9, 42, 4, 43, 9, 43, 4, 44, 9, 44,
	4, 45, 9, 45, 4, 46, 9, 46, 4, 47, 9, 47, 4, 48, 9, 48, 4, 49, 9, 49, 4,
	50, 9, 50, 4, 51, 9, 51, 4, 52, 9, 52, 4, 53, 9, 53, 4, 54, 9, 54, 4, 55,
	9, 55, 4, 56, 9, 56, 4, 57, 9, 57, 4, 58, 9, 58, 4, 59, 9, 59, 4, 60, 9,
	60, 4, 61, 9, 61, 4, 62, 9, 62, 4, 63, 9, 63, 4, 64, 9, 64, 4, 65, 9, 65,
	4, 66, 9, 66, 4, 67, 9, 67, 4, 68, 9, 68, 4, 69, 9, 69, 4, 70, 9, 70, 3,
	2, 7, 2, 142, 10, 2, 12, 2, 14, 2, 145, 11, 2, 3, 2, 3, 2, 3, 3, 3, 3,
	3, 4, 3, 4, 3, 5, 3, 5, 3, 5, 3, 6, 7, 6, 157, 10, 6, 12, 6, 14, 6, 160,
	11, 6, 3, 6, 3, 6, 3, 7, 3, 7, 3, 7, 5, 7, 167, 10, 7, 3, 8, 3, 8, 5, 8,
	171, 10, 8, 3, 9, 3, 9, 5, 9, 175, 10, 9, 3, 9, 3, 9, 3, 9, 5, 9, 180,
	10, 9, 3, 9, 3, 9, 3, 9, 3, 9, 3, 9, 3, 9, 5, 9, 188, 10, 9, 3, 10, 3,
	10, 3, 10, 3, 10, 5, 10, 194, 10, 10, 3, 10, 3, 10, 3, 10, 7, 10, 199,
	10, 10, 12, 10, 14, 10, 202, 11, 10, 3, 10, 3, 10, 3, 10, 3, 10, 3, 10,
	5, 10, 209, 10, 10, 3, 10, 3, 10, 3, 10, 7, 10, 214, 10, 10, 12, 10, 14,
	10, 217, 11, 10, 3, 10, 3, 10, 5, 10, 221, 10, 10, 3, 11, 3, 11, 3, 12,
	3, 12, 3, 13, 3, 13, 3, 13, 3, 13, 3, 13, 3, 13, 3, 13, 5, 13, 234, 10,
	13, 3, 14, 3, 14, 3, 14, 3, 14, 5, 14, 240, 10, 14, 3, 15, 3, 15, 5, 15,
	244, 10, 15, 3, 16, 3, 16, 5, 16, 248, 10, 16, 3, 17, 3, 17, 5, 17, 252,
	10, 17, 3, 18, 3, 18, 3, 18, 3, 19, 3, 19, 3, 19, 3, 19, 5, 19, 261, 10,
	19, 3, 20, 3, 20, 5, 20, 265, 10, 20, 3, 21, 3, 21, 3, 21, 3, 21, 7, 21,
	271, 10, 21, 12, 21, 14, 21, 274, 11, 21, 3, 22, 3, 22, 5, 22, 278, 10,
	22, 3, 23, 3, 23, 3, 23, 3, 23, 3, 23, 3, 23, 3, 23, 3, 23, 3, 23, 3, 23,
	3, 23, 3, 23, 3, 23, 3, 23, 3, 23, 3, 23, 3, 23, 3, 23, 5, 23, 298, 10,
	23, 3, 24, 3, 24, 3, 24, 3, 24, 3, 25, 3, 25, 3, 25, 7, 25, 307, 10, 25,
	12, 25, 14, 25, 310, 11, 25, 3, 26, 3, 26, 3, 26, 3, 26, 7, 26, 316, 10,
	26, 12, 26, 14, 26, 319, 11, 26, 3, 27, 3, 27, 3, 27, 3, 27, 3, 28, 3,
	28, 3, 28, 3, 28, 3, 28, 3, 28, 5, 28, 331, 10, 28, 5, 28, 333, 10, 28,
	3, 29, 3, 29, 3, 29, 3, 29, 3, 29, 3, 30, 3, 30, 3, 31, 3, 31, 3, 31, 3,
	31, 5, 31, 346, 10, 31, 3, 32, 3, 32, 3, 32, 3, 32, 3, 32, 5, 32, 353,
	10, 32, 3, 33, 3, 33, 3, 33, 5, 33, 358, 10, 33, 3, 34, 3, 34, 3, 34, 3,
	34, 3, 34, 3, 34, 5, 34, 366, 10, 34, 3, 35, 3, 35, 3, 35, 3, 35, 3, 35,
	3, 35, 3, 35, 3, 35, 3, 35, 3, 35, 3, 35, 3, 35, 3, 35, 3, 35, 3, 35, 5,
	35, 383, 10, 35, 3, 36, 3, 36, 3, 36, 3, 37, 3, 37, 3, 38, 3, 38, 3, 38,
	5, 38, 393, 10, 38, 3, 38, 3, 38, 3, 38, 3, 38, 5, 38, 399, 10, 38, 3,
	39, 3, 39, 5, 39, 403, 10, 39, 3, 39, 3, 39, 3, 40, 3, 40, 3, 40, 3, 40,
	7, 40, 411, 10, 40, 12, 40, 14, 40, 414, 11, 40, 5, 40, 416, 10, 40, 3,
	40, 5, 40, 419, 10, 40, 3, 40, 3, 40, 3, 41, 3, 41, 3, 42, 3, 42, 3, 43,
	3, 43, 3, 44, 3, 44, 3, 45, 3, 45, 3, 46, 3, 46, 6, 46, 435, 10, 46, 13,
	46, 14, 46, 436, 3, 46, 7, 46, 440, 10, 46, 12, 46, 14, 46, 443, 11, 46,
	3, 47, 3, 47, 3, 47, 3, 47, 3, 47, 3, 47, 3, 47, 3, 47, 3, 47, 5, 47, 454,
	10, 47, 3, 48, 3, 48, 3, 49, 3, 49, 3, 49, 3, 49, 3, 50, 3, 50, 3, 50,
	5, 50, 465, 10, 50, 3, 51, 3, 51, 3, 51, 3, 51, 3, 52, 3, 52, 3, 52, 3,
	53, 7, 53, 475, 10, 53, 12, 53, 14, 53, 478, 11, 53, 3, 54, 3, 54, 3, 55,
	3, 55, 3, 55, 3, 55, 3, 56, 3, 56, 3, 56, 5, 56, 489, 10, 56, 3, 57, 3,
	57, 3, 57, 7, 57, 494, 10, 57, 12, 57, 14, 57, 497, 11, 57, 6, 57, 499,
	10, 57, 13, 57, 14, 57, 500, 3, 57, 3, 57, 3, 57, 3, 57, 7, 57, 507, 10,
	57, 12, 57, 14, 57, 510, 11, 57, 7, 57, 512, 10, 57, 12, 57, 14, 57, 515,
	11, 57, 3, 57, 3, 57, 3, 57, 7, 57, 520, 10, 57, 12, 57, 14, 57, 523, 11,
	57, 7, 57, 525, 10, 57, 12, 57, 14, 57, 528, 11, 57, 5, 57, 530, 10, 57,
	3, 58, 3, 58, 3, 58, 3, 59, 3, 59, 3, 59, 3, 59, 7, 59, 539, 10, 59, 12,
	59, 14, 59, 542, 11, 59, 5, 59, 544, 10, 59, 3, 59, 3, 59, 3, 60, 3, 60,
	3, 60, 3, 60, 3, 60, 3, 60, 3, 60, 3, 60, 3, 60, 3, 60, 3, 60, 3, 60, 3,
	60, 3, 60, 3, 60, 3, 60, 3, 60, 5, 60, 565, 10, 60, 3, 60, 3, 60, 3, 60,
	3, 60, 3, 60, 3, 60, 3, 60, 3, 60, 3, 60, 3, 60, 3, 60, 3, 60, 5, 60, 579,
	10, 60, 3, 60, 3, 60, 3, 60, 3, 60, 3, 60, 3, 60, 3, 60, 3, 60, 3, 60,
	3, 60, 3, 60, 3, 60, 3, 60, 3, 60, 3, 60, 3, 60, 3, 60, 3, 60, 3, 60, 3,
	60, 3, 60, 3, 60, 3, 60, 3, 60, 3, 60, 5, 60, 606, 10, 60, 3, 60, 3, 60,
	7, 60, 610, 10, 60, 12, 60, 14, 60, 613, 11, 60, 3, 61, 3, 61, 3, 61, 5,
	61, 618, 10, 61, 3, 61, 3, 61, 3, 61, 3, 61, 3, 61, 3, 61, 3, 61, 3, 61,
	3, 61, 3, 61, 3, 61, 3, 61, 3, 61, 3, 61, 3, 61, 3, 61, 3, 61, 3, 61, 3,
	61, 3, 61, 3, 61, 3, 61, 3, 61, 5, 61, 643, 10, 61, 3, 62, 3, 62, 3, 63,
	3, 63, 3, 63, 5, 63, 650, 10, 63, 3, 64, 3, 64, 3, 65, 3, 65, 3, 66, 3,
	66, 3, 67, 3, 67, 3, 68, 3, 68, 3, 69, 3, 69, 3, 70, 3, 70, 3, 70, 2, 3,
	118, 71, 2, 4, 6, 8, 10, 12, 14, 16, 18, 20, 22, 24, 26, 28, 30, 32, 34,
	36, 38, 40, 42, 44, 46, 48, 50, 52, 54, 56, 58, 60, 62, 64, 66, 68, 70,
	72, 74, 76, 78, 80, 82, 84, 86, 88, 90, 92, 94, 96, 98, 100, 102, 104,
	106, 108, 110, 112, 114, 116, 118, 120, 122, 124, 126, 128, 130, 132, 134,
	136, 138, 2, 10, 3, 2, 47, 48, 5, 2, 30, 31, 37, 61, 65, 65, 4, 2, 47,
	47, 55, 56, 3, 2, 17, 22, 3, 2, 35, 36, 3, 2, 23, 25, 3, 2, 26, 27, 4,
	2, 26, 27, 59, 60, 2, 701, 2, 143, 3, 2, 2, 2, 4, 148, 3, 2, 2, 2, 6, 150,
	3, 2, 2, 2, 8, 152, 3, 2, 2, 2, 10, 158, 3, 2, 2, 2, 12, 166, 3, 2, 2,
	2, 14, 170, 3, 2, 2, 2, 16, 187, 3, 2, 2, 2, 18, 220, 3, 2, 2, 2, 20, 222,
	3, 2, 2, 2, 22, 224, 3, 2, 2, 2, 24, 233, 3, 2, 2, 2, 26, 239, 3, 2, 2,
	2, 28, 243, 3, 2, 2, 2, 30, 247, 3, 2, 2, 2, 32, 251, 3, 2, 2, 2, 34, 253,
	3, 2, 2, 2, 36, 256, 3, 2, 2, 2, 38, 264, 3, 2, 2, 2, 40, 266, 3, 2, 2,
	2, 42, 275, 3, 2, 2, 2, 44, 297, 3, 2, 2, 2, 46, 299, 3, 2, 2, 2, 48, 303,
	3, 2, 2, 2, 50, 311, 3, 2, 2, 2, 52, 320, 3, 2, 2, 2, 54, 332, 3, 2, 2,
	2, 56, 334, 3, 2, 2, 2, 58, 339, 3, 2, 2, 2, 60, 345, 3, 2, 2, 2, 62, 352,
	3, 2, 2, 2, 64, 357, 3, 2, 2, 2, 66, 359, 3, 2, 2, 2, 68, 382, 3, 2, 2,
	2, 70, 384, 3, 2, 2, 2, 72, 387, 3, 2, 2, 2, 74, 392, 3, 2, 2, 2, 76, 400,
	3, 2, 2, 2, 78, 406, 3, 2, 2, 2, 80, 422, 3, 2, 2, 2, 82, 424, 3, 2, 2,
	2, 84, 426, 3, 2, 2, 2, 86, 428, 3, 2, 2, 2, 88, 430, 3, 2, 2, 2, 90, 432,
	3, 2, 2, 2, 92, 453, 3, 2, 2, 2, 94, 455, 3, 2, 2, 2, 96, 457, 3, 2, 2,
	2, 98, 464, 3, 2, 2, 2, 100, 466, 3, 2, 2, 2, 102, 470, 3, 2, 2, 2, 104,
	476, 3, 2, 2, 2, 106, 479, 3, 2, 2, 2, 108, 481, 3, 2, 2, 2, 110, 488,
	3, 2, 2, 2, 112, 529, 3, 2, 2, 2, 114, 531, 3, 2, 2, 2, 116, 534, 3, 2,
	2, 2, 118, 564, 3, 2, 2, 2, 120, 642, 3, 2, 2, 2, 122, 644, 3, 2, 2, 2,
	124, 649, 3, 2, 2, 2, 126, 651, 3, 2, 2, 2, 128, 653, 3, 2, 2, 2, 130,
	655, 3, 2, 2, 2, 132, 657, 3, 2, 2, 2, 134, 659, 3, 2, 2, 2, 136, 661,
	3, 2, 2, 2, 138, 663, 3, 2, 2, 2, 140, 142, 5, 4, 3, 2, 141, 140, 3, 2,
	2, 2, 142, 145, 3, 2, 2, 2, 143, 141, 3, 2, 2, 2, 143, 144, 3, 2, 2, 2,
	144, 146, 3, 2, 2, 2, 145, 143, 3, 2, 2, 2, 146, 147, 5, 10, 6, 2, 147,
	3, 3, 2, 2, 2, 148, 149, 5, 6, 4, 2, 149, 5, 3, 2, 2, 2, 150, 151, 5, 8,
	5, 2, 151, 7, 3, 2, 2, 2, 152, 153, 7, 50, 2, 2, 153, 154, 5, 102, 52,
	2, 154, 9, 3, 2, 2, 2, 155, 157, 5, 12, 7, 2, 156, 155, 3, 2, 2, 2, 157,
	160, 3, 2, 2, 2, 158, 156, 3, 2, 2, 2, 158, 159, 3, 2, 2, 2, 159, 161,
	3, 2, 2, 2, 160, 158, 3, 2, 2, 2, 161, 162, 5, 14, 8, 2, 162, 11, 3, 2,
	2, 2, 163, 167, 5, 108, 55, 2, 164, 167, 5, 68, 35, 2, 165, 167, 5, 58,
	30, 2, 166, 163, 3, 2, 2, 2, 166, 164, 3, 2, 2, 2, 166, 165, 3, 2, 2, 2,
	167, 13, 3, 2, 2, 2, 168, 171, 5, 16, 9, 2, 169, 171, 5, 18, 10, 2, 170,
	168, 3, 2, 2, 2, 170, 169, 3, 2, 2, 2, 171, 15, 3, 2, 2, 2, 172, 174, 7,
	38, 2, 2, 173, 175, 7, 40, 2, 2, 174, 173, 3, 2, 2, 2, 174, 175, 3, 2,
	2, 2, 175, 176, 3, 2, 2, 2, 176, 188, 5, 118, 60, 2, 177, 179, 7, 38, 2,
	2, 178, 180, 7, 40, 2, 2, 179, 178, 3, 2, 2, 2, 179, 180, 3, 2, 2, 2, 180,
	181, 3, 2, 2, 2, 181, 182, 7, 13, 2, 2, 182, 183, 5, 18, 10, 2, 183, 184,
	7, 14, 2, 2, 184, 188, 3, 2, 2, 2, 185, 186, 7, 38, 2, 2, 186, 188, 5,
	120, 61, 2, 187, 172, 3, 2, 2, 2, 187, 177, 3, 2, 2, 2, 187, 185, 3, 2,
	2, 2, 188, 17, 3, 2, 2, 2, 189, 190, 7, 37, 2, 2, 190, 193, 5, 20, 11,
	2, 191, 192, 7, 10, 2, 2, 192, 194, 5, 22, 12, 2, 193, 191, 3, 2, 2, 2,
	193, 194, 3, 2, 2, 2, 194, 195, 3, 2, 2, 2, 195, 196, 7, 61, 2, 2, 196,
	200, 5, 24, 13, 2, 197, 199, 5, 30, 16, 2, 198, 197, 3, 2, 2, 2, 199, 202,
	3, 2, 2, 2, 200, 198, 3, 2, 2, 2, 200, 201, 3, 2, 2, 2, 201, 203, 3, 2,
	2, 2, 202, 200, 3, 2, 2, 2, 203, 204, 5, 32, 17, 2, 204, 221, 3, 2, 2,
	2, 205, 206, 7, 37, 2, 2, 206, 208, 5, 20, 11, 2, 207, 209, 7, 62, 2, 2,
	208, 207, 3, 2, 2, 2, 208, 209, 3, 2, 2, 2, 209, 210, 3, 2, 2, 2, 210,
	211, 7, 63, 2, 2, 211, 215, 5, 118, 60, 2, 212, 214, 5, 30, 16, 2, 213,
	212, 3, 2, 2, 2, 214, 217, 3, 2, 2, 2, 215, 213, 3, 2, 2, 2, 215, 216,
	3, 2, 2, 2, 216, 218, 3, 2, 2, 2, 217, 215, 3, 2, 2, 2, 218, 219, 5, 32,
	17, 2, 219, 221, 3, 2, 2, 2, 220, 189, 3, 2, 2, 2, 220, 205, 3, 2, 2, 2,
	221, 19, 3, 2, 2, 2, 222, 223, 7, 65, 2, 2, 223, 21, 3, 2, 2, 2, 224, 225,
	7, 65, 2, 2, 225, 23, 3, 2, 2, 2, 226, 234, 5, 108, 55, 2, 227, 234, 5,
	76, 39, 2, 228, 234, 5, 78, 40, 2, 229, 234, 5, 72, 37, 2, 230, 234, 5,
	114, 58, 2, 231, 234, 5, 74, 38, 2, 232, 234, 5, 70, 36, 2, 233, 226, 3,
	2, 2, 2, 233, 227, 3, 2, 2, 2, 233, 228, 3, 2, 2, 2, 233, 229, 3, 2, 2,
	2, 233, 230, 3, 2, 2, 2, 233, 231, 3, 2, 2, 2, 233, 232, 3, 2, 2, 2, 234,
	25, 3, 2, 2, 2, 235, 240, 5, 36, 19, 2, 236, 240, 5, 40, 21, 2, 237, 240,
	5, 34, 18, 2, 238, 240, 5, 44, 23, 2, 239, 235, 3, 2, 2, 2, 239, 236, 3,
	2, 2, 2, 239, 237, 3, 2, 2, 2, 239, 238, 3, 2, 2, 2, 240, 27, 3, 2, 2,
	2, 241, 244, 5, 68, 35, 2, 242, 244, 5, 108, 55, 2, 243, 241, 3, 2, 2,
	2, 243, 242, 3, 2, 2, 2, 244, 29, 3, 2, 2, 2, 245, 248, 5, 28, 15, 2, 246,
	248, 5, 26, 14, 2, 247, 245, 3, 2, 2, 2, 247, 246, 3, 2, 2, 2, 248, 31,
	3, 2, 2, 2, 249, 252, 5, 16, 9, 2, 250, 252, 5, 18, 10, 2, 251, 249, 3,
	2, 2, 2, 251, 250, 3, 2, 2, 2, 252, 33, 3, 2, 2, 2, 253, 254, 7, 41, 2,
	2, 254, 255, 5, 118, 60, 2, 255, 35, 3, 2, 2, 2, 256, 257, 7, 43, 2, 2,
	257, 260, 5, 38, 20, 2, 258, 259, 7, 10, 2, 2, 259, 261, 5, 38, 20, 2,
	260, 258, 3, 2, 2, 2, 260, 261, 3, 2, 2, 2, 261, 37, 3, 2, 2, 2, 262, 265,
	7, 67, 2, 2, 263, 265, 5, 70, 36, 2, 264, 262, 3, 2, 2, 2, 264, 263, 3,
	2, 2, 2, 265, 39, 3, 2, 2, 2, 266, 267, 7, 42, 2, 2, 267, 272, 5, 42, 22,
	2, 268, 269, 7, 10, 2, 2, 269, 271, 5, 42, 22, 2, 270, 268, 3, 2, 2, 2,
	271, 274, 3, 2, 2, 2, 272, 270, 3, 2, 2, 2, 272, 273, 3, 2, 2, 2, 273,
	41, 3, 2, 2, 2, 274, 272, 3, 2, 2, 2, 275, 277, 5, 118, 60, 2, 276, 278,
	7, 46, 2, 2, 277, 276, 3, 2, 2, 2, 277, 278, 3, 2, 2, 2, 278, 43, 3, 2,
	2, 2, 279, 280, 7, 45, 2, 2, 280, 298, 5, 56, 29, 2, 281, 282, 7, 45, 2,
	2, 282, 298, 5, 50, 26, 2, 283, 284, 7, 45, 2, 2, 284, 285, 5, 48, 25,
	2, 285, 286, 5, 50, 26, 2, 286, 298, 3, 2, 2, 2, 287, 288, 7, 45, 2, 2,
	288, 289, 5, 48, 25, 2, 289, 290, 5, 54, 28, 2, 290, 298, 3, 2, 2, 2, 291,
	292, 7, 45, 2, 2, 292, 293, 5, 48, 25, 2, 293, 294, 5, 56, 29, 2, 294,
	298, 3, 2, 2, 2, 295, 296, 7, 45, 2, 2, 296, 298, 5, 48, 25, 2, 297, 279,
	3, 2, 2, 2, 297, 281, 3, 2, 2, 2, 297, 283, 3, 2, 2, 2, 297, 287, 3, 2,
	2, 2, 297, 291, 3, 2, 2, 2, 297, 295, 3, 2, 2, 2, 298, 45, 3, 2, 2, 2,
	299, 300, 7, 65, 2, 2, 300, 301, 7, 33, 2, 2, 301, 302, 5, 118, 60, 2,
	302, 47, 3, 2, 2, 2, 303, 308, 5, 46, 24, 2, 304, 305, 7, 10, 2, 2, 305,
	307, 5, 46, 24, 2, 306, 304, 3, 2, 2, 2, 307, 310, 3, 2, 2, 2, 308, 306,
	3, 2, 2, 2, 308, 309, 3, 2, 2, 2, 309, 49, 3, 2, 2, 2, 310, 308, 3, 2,
	2, 2, 311, 312, 7, 57, 2, 2, 312, 317, 5, 52, 27, 2, 313, 314, 7, 10, 2,
	2, 314, 316, 5, 52, 27, 2, 315, 313, 3, 2, 2, 2, 316, 319, 3, 2, 2, 2,
	317, 315, 3, 2, 2, 2, 317, 318, 3, 2, 2, 2, 318, 51, 3, 2, 2, 2, 319, 317,
	3, 2, 2, 2, 320, 321, 7, 65, 2, 2, 321, 322, 7, 33, 2, 2, 322, 323, 5,
	108, 55, 2, 323, 53, 3, 2, 2, 2, 324, 325, 7, 51, 2, 2, 325, 333, 5, 46,
	24, 2, 326, 327, 7, 51, 2, 2, 327, 330, 7, 65, 2, 2, 328, 329, 7, 52, 2,
	2, 329, 331, 7, 65, 2, 2, 330, 328, 3, 2, 2, 2, 330, 331, 3, 2, 2, 2, 331,
	333, 3, 2, 2, 2, 332, 324, 3, 2, 2, 2, 332, 326, 3, 2, 2, 2, 333, 55, 3,
	2, 2, 2, 334, 335, 7, 53, 2, 2, 335, 336, 7, 54, 2, 2, 336, 337, 7, 51,
	2, 2, 337, 338, 7, 65, 2, 2, 338, 57, 3, 2, 2, 2, 339, 340, 5, 66, 34,
	2, 340, 59, 3, 2, 2, 2, 341, 346, 5, 84, 43, 2, 342, 346, 5, 72, 37, 2,
	343, 346, 5, 108, 55, 2, 344, 346, 5, 114, 58, 2, 345, 341, 3, 2, 2, 2,
	345, 342, 3, 2, 2, 2, 345, 343, 3, 2, 2, 2, 345, 344, 3, 2, 2, 2, 346,
	61, 3, 2, 2, 2, 347, 353, 5, 82, 42, 2, 348, 353, 5, 72, 37, 2, 349, 353,
	5, 70, 36, 2, 350, 353, 5, 108, 55, 2, 351, 353, 5, 114, 58, 2, 352, 347,
	3, 2, 2, 2, 352, 348, 3, 2, 2, 2, 352, 349, 3, 2, 2, 2, 352, 350, 3, 2,
	2, 2, 352, 351, 3, 2, 2, 2, 353, 63, 3, 2, 2, 2, 354, 358, 5, 72, 37, 2,
	355, 358, 5, 108, 55, 2, 356, 358, 5, 114, 58, 2, 357, 354, 3, 2, 2, 2,
	357, 355, 3, 2, 2, 2, 357, 356, 3, 2, 2, 2, 358, 65, 3, 2, 2, 2, 359, 360,
	7, 39, 2, 2, 360, 361, 7, 58, 2, 2, 361, 362, 5, 62, 32, 2, 362, 363, 7,
	61, 2, 2, 363, 365, 5, 64, 33, 2, 364, 366, 5, 60, 31, 2, 365, 364, 3,
	2, 2, 2, 365, 366, 3, 2, 2, 2, 366, 67, 3, 2, 2, 2, 367, 368, 7, 44, 2,
	2, 368, 369, 7, 65, 2, 2, 369, 370, 7, 33, 2, 2, 370, 383, 5, 118, 60,
	2, 371, 372, 7, 44, 2, 2, 372, 373, 7, 65, 2, 2, 373, 374, 7, 33, 2, 2,
	374, 375, 7, 13, 2, 2, 375, 376, 5, 18, 10, 2, 376, 377, 7, 14, 2, 2, 377,
	383, 3, 2, 2, 2, 378, 379, 7, 44, 2, 2, 379, 380, 7, 65, 2, 2, 380, 381,
	7, 33, 2, 2, 381, 383, 5, 120, 61, 2, 382, 367, 3, 2, 2, 2, 382, 371, 3,
	2, 2, 2, 382, 378, 3, 2, 2, 2, 383, 69, 3, 2, 2, 2, 384, 385, 7, 64, 2,
	2, 385, 386, 7, 65, 2, 2, 386, 71, 3, 2, 2, 2, 387, 388, 7, 65, 2, 2, 388,
	73, 3, 2, 2, 2, 389, 393, 5, 84, 43, 2, 390, 393, 5, 72, 37, 2, 391, 393,
	5, 70, 36, 2, 392, 389, 3, 2, 2, 2, 392, 390, 3, 2, 2, 2, 392, 391, 3,
	2, 2, 2, 393, 394, 3, 2, 2, 2, 394, 398, 7, 32, 2, 2, 395, 399, 5, 84,
	43, 2, 396, 399, 5, 72, 37, 2, 397, 399, 5, 70, 36, 2, 398, 395, 3, 2,
	2, 2, 398, 396, 3, 2, 2, 2, 398, 397, 3, 2, 2, 2, 399, 75, 3, 2, 2, 2,
	400, 402, 7, 11, 2, 2, 401, 403, 5, 90, 46, 2, 402, 401, 3, 2, 2, 2, 402,
	403, 3, 2, 2, 2, 403, 404, 3, 2, 2, 2, 404, 405, 7, 12, 2, 2, 405, 77,
	3, 2, 2, 2, 406, 415, 7, 15, 2, 2, 407, 412, 5, 92, 47, 2, 408, 409, 7,
	10, 2, 2, 409, 411, 5, 92, 47, 2, 410, 408, 3, 2, 2, 2, 411, 414, 3, 2,
	2, 2, 412, 410, 3, 2, 2, 2, 412, 413, 3, 2, 2, 2, 413, 416, 3, 2, 2, 2,
	414, 412, 3, 2, 2, 2, 415, 407, 3, 2, 2, 2, 415, 416, 3, 2, 2, 2, 416,
	418, 3, 2, 2, 2, 417, 419, 7, 10, 2, 2, 418, 417, 3, 2, 2, 2, 418, 419,
	3, 2, 2, 2, 419, 420, 3, 2, 2, 2, 420, 421, 7, 16, 2, 2, 421, 79, 3, 2,
	2, 2, 422, 423, 7, 49, 2, 2, 423, 81, 3, 2, 2, 2, 424, 425, 7, 66, 2, 2,
	425, 83, 3, 2, 2, 2, 426, 427, 7, 67, 2, 2, 427, 85, 3, 2, 2, 2, 428, 429,
	7, 68, 2, 2, 429, 87, 3, 2, 2, 2, 430, 431, 9, 2, 2, 2, 431, 89, 3, 2,
	2, 2, 432, 441, 5, 118, 60, 2, 433, 435, 7, 10, 2, 2, 434, 433, 3, 2, 2,
	2, 435, 436, 3, 2, 2, 2, 436, 434, 3, 2, 2, 2, 436, 437, 3, 2, 2, 2, 437,
	438, 3, 2, 2, 2, 438, 440, 5, 118, 60, 2, 439, 434, 3, 2, 2, 2, 440, 443,
	3, 2, 2, 2, 441, 439, 3, 2, 2, 2, 441, 442, 3, 2, 2, 2, 442, 91, 3, 2,
	2, 2, 443, 441, 3, 2, 2, 2, 444, 445, 5, 98, 50, 2, 445, 446, 7, 7, 2,
	2, 446, 447, 5, 118, 60, 2, 447, 454, 3, 2, 2, 2, 448, 449, 5, 96, 49,
	2, 449, 450, 7, 7, 2, 2, 450, 451, 5, 118, 60, 2, 451, 454, 3, 2, 2, 2,
	452, 454, 5, 94, 48, 2, 453, 444, 3, 2, 2, 2, 453, 448, 3, 2, 2, 2, 453,
	452, 3, 2, 2, 2, 454, 93, 3, 2, 2, 2, 455, 456, 5, 72, 37, 2, 456, 95,
	3, 2, 2, 2, 457, 458, 7, 11, 2, 2, 458, 459, 5, 118, 60, 2, 459, 460, 7,
	12, 2, 2, 460, 97, 3, 2, 2, 2, 461, 465, 7, 65, 2, 2, 462, 465, 5, 82,
	42, 2, 463, 465, 5, 70, 36, 2, 464, 461, 3, 2, 2, 2, 464, 462, 3, 2, 2,
	2, 464, 463, 3, 2, 2, 2, 465, 99, 3, 2, 2, 2, 466, 467, 7, 13, 2, 2, 467,
	468, 5, 118, 60, 2, 468, 469, 7, 14, 2, 2, 469, 101, 3, 2, 2, 2, 470, 471,
	5, 104, 53, 2, 471, 472, 7, 65, 2, 2, 472, 103, 3, 2, 2, 2, 473, 475, 7,
	69, 2, 2, 474, 473, 3, 2, 2, 2, 475, 478, 3, 2, 2, 2, 476, 474, 3, 2, 2,
	2, 476, 477, 3, 2, 2, 2, 477, 105, 3, 2, 2, 2, 478, 476, 3, 2, 2, 2, 479,
	480, 9, 3, 2, 2, 480, 107, 3, 2, 2, 2, 481, 482, 5, 104, 53, 2, 482, 483,
	5, 106, 54, 2, 483, 484, 5, 116, 59, 2, 484, 109, 3, 2, 2, 2, 485, 489,
	7, 65, 2, 2, 486, 489, 5, 108, 55, 2, 487, 489, 5, 70, 36, 2, 488, 485,
	3, 2, 2, 2, 488, 486, 3, 2, 2, 2, 488, 487, 3, 2, 2, 2, 489, 111, 3, 2,
	2, 2, 490, 491, 7, 9, 2, 2, 491, 495, 5, 98, 50, 2, 492, 494, 5, 96, 49,
	2, 493, 492, 3, 2, 2, 2, 494, 497, 3, 2, 2, 2, 495, 493, 3, 2, 2, 2, 495,
	496, 3, 2, 2, 2, 496, 499, 3, 2, 2, 2, 497, 495, 3, 2, 2, 2, 498, 490,
	3, 2, 2, 2, 499, 500, 3, 2, 2, 2, 500, 498, 3, 2, 2, 2, 500, 501, 3, 2,
	2, 2, 501, 530, 3, 2, 2, 2, 502, 513, 5, 96, 49, 2, 503, 504, 7, 9, 2,
	2, 504, 508, 5, 98, 50, 2, 505, 507, 5, 96, 49, 2, 506, 505, 3, 2, 2, 2,
	507, 510, 3, 2, 2, 2, 508, 506, 3, 2, 2, 2, 508, 509, 3, 2, 2, 2, 509,
	512, 3, 2, 2, 2, 510, 508, 3, 2, 2, 2, 511, 503, 3, 2, 2, 2, 512, 515,
	3, 2, 2, 2, 513, 511, 3, 2, 2, 2, 513, 514, 3, 2, 2, 2, 514, 526, 3, 2,
	2, 2, 515, 513, 3, 2, 2, 2, 516, 521, 5, 96, 49, 2, 517, 518, 7, 9, 2,
	2, 518, 520, 5, 98, 50, 2, 519, 517, 3, 2, 2, 2, 520, 523, 3, 2, 2, 2,
	521, 519, 3, 2, 2, 2, 521, 522, 3, 2, 2, 2, 522, 525, 3, 2, 2, 2, 523,
	521, 3, 2, 2, 2, 524, 516, 3, 2, 2, 2, 525, 528, 3, 2, 2, 2, 526, 524,
	3, 2, 2, 2, 526, 527, 3, 2, 2, 2, 527, 530, 3, 2, 2, 2, 528, 526, 3, 2,
	2, 2, 529, 498, 3, 2, 2, 2, 529, 502, 3, 2, 2, 2, 530, 113, 3, 2, 2, 2,
	531, 532, 5, 110, 56, 2, 532, 533, 5, 112, 57, 2, 533, 115, 3, 2, 2, 2,
	534, 543, 7, 13, 2, 2, 535, 540, 5, 118, 60, 2, 536, 537, 7, 10, 2, 2,
	537, 539, 5, 118, 60, 2, 538, 536, 3, 2, 2, 2, 539, 542, 3, 2, 2, 2, 540,
	538, 3, 2, 2, 2, 540, 541, 3, 2, 2, 2, 541, 544, 3, 2, 2, 2, 542, 540,
	3, 2, 2, 2, 543, 535, 3, 2, 2, 2, 543, 544, 3, 2, 2, 2, 544, 545, 3, 2,
	2, 2, 545, 546, 7, 14, 2, 2, 546, 117, 3, 2, 2, 2, 547, 548, 8, 60, 1,
	2, 548, 549, 5, 138, 70, 2, 549, 550, 5, 118, 60, 25, 550, 565, 3, 2, 2,
	2, 551, 565, 5, 108, 55, 2, 552, 565, 5, 100, 51, 2, 553, 565, 5, 74, 38,
	2, 554, 565, 5, 82, 42, 2, 555, 565, 5, 84, 43, 2, 556, 565, 5, 86, 44,
	2, 557, 565, 5, 80, 41, 2, 558, 565, 5, 76, 39, 2, 559, 565, 5, 78, 40,
	2, 560, 565, 5, 72, 37, 2, 561, 565, 5, 114, 58, 2, 562, 565, 5, 88, 45,
	2, 563, 565, 5, 70, 36, 2, 564, 547, 3, 2, 2, 2, 564, 551, 3, 2, 2, 2,
	564, 552, 3, 2, 2, 2, 564, 553, 3, 2, 2, 2, 564, 554, 3, 2, 2, 2, 564,
	555, 3, 2, 2, 2, 564, 556, 3, 2, 2, 2, 564, 557, 3, 2, 2, 2, 564, 558,
	3, 2, 2, 2, 564, 559, 3, 2, 2, 2, 564, 560, 3, 2, 2, 2, 564, 561, 3, 2,
	2, 2, 564, 562, 3, 2, 2, 2, 564, 563, 3, 2, 2, 2, 565, 611, 3, 2, 2, 2,
	566, 567, 12, 24, 2, 2, 567, 568, 5, 134, 68, 2, 568, 569, 5, 118, 60,
	25, 569, 610, 3, 2, 2, 2, 570, 571, 12, 23, 2, 2, 571, 572, 5, 136, 69,
	2, 572, 573, 5, 118, 60, 24, 573, 610, 3, 2, 2, 2, 574, 575, 12, 20, 2,
	2, 575, 578, 5, 122, 62, 2, 576, 579, 5, 124, 63, 2, 577, 579, 5, 126,
	64, 2, 578, 576, 3, 2, 2, 2, 578, 577, 3, 2, 2, 2, 579, 580, 3, 2, 2, 2,
	580, 581, 5, 118, 60, 21, 581, 610, 3, 2, 2, 2, 582, 583, 12, 19, 2, 2,
	583, 584, 5, 124, 63, 2, 584, 585, 5, 118, 60, 20, 585, 610, 3, 2, 2, 2,
	586, 587, 12, 18, 2, 2, 587, 588, 5, 126, 64, 2, 588, 589, 5, 118, 60,
	19, 589, 610, 3, 2, 2, 2, 590, 591, 12, 17, 2, 2, 591, 592, 5, 128, 65,
	2, 592, 593, 5, 118, 60, 18, 593, 610, 3, 2, 2, 2, 594, 595, 12, 16, 2,
	2, 595, 596, 5, 130, 66, 2, 596, 597, 5, 118, 60, 17, 597, 610, 3, 2, 2,
	2, 598, 599, 12, 15, 2, 2, 599, 600, 5, 132, 67, 2, 600, 601, 5, 118, 60,
	16, 601, 610, 3, 2, 2, 2, 602, 603, 12, 14, 2, 2, 603, 605, 7, 34, 2, 2,
	604, 606, 5, 118, 60, 2, 605, 604, 3, 2, 2, 2, 605, 606, 3, 2, 2, 2, 606,
	607, 3, 2, 2, 2, 607, 608, 7, 7, 2, 2, 608, 610, 5, 118, 60, 15, 609, 566,
	3, 2, 2, 2, 609, 570, 3, 2, 2, 2, 609, 574, 3, 2, 2, 2, 609, 582, 3, 2,
	2, 2, 609, 586, 3, 2, 2, 2, 609, 590, 3, 2, 2, 2, 609, 594, 3, 2, 2, 2,
	609, 598, 3, 2, 2, 2, 609, 602, 3, 2, 2, 2, 610, 613, 3, 2, 2, 2, 611,
	609, 3, 2, 2, 2, 611, 612, 3, 2, 2, 2, 612, 119, 3, 2, 2, 2, 613, 611,
	3, 2, 2, 2, 614, 615, 5, 118, 60, 2, 615, 617, 7, 34, 2, 2, 616, 618, 5,
	118, 60, 2, 617, 616, 3, 2, 2, 2, 617, 618, 3, 2, 2, 2, 618, 619, 3, 2,
	2, 2, 619, 620, 7, 7, 2, 2, 620, 621, 7, 13, 2, 2, 621, 622, 5, 18, 10,
	2, 622, 623, 7, 14, 2, 2, 623, 643, 3, 2, 2, 2, 624, 625, 5, 118, 60, 2,
	625, 626, 7, 34, 2, 2, 626, 627, 7, 13, 2, 2, 627, 628, 5, 18, 10, 2, 628,
	629, 7, 14, 2, 2, 629, 630, 7, 7, 2, 2, 630, 631, 5, 118, 60, 2, 631, 643,
	3, 2, 2, 2, 632, 633, 5, 118, 60, 2, 633, 634, 7, 34, 2, 2, 634, 635, 7,
	13, 2, 2, 635, 636, 5, 18, 10, 2, 636, 637, 7, 14, 2, 2, 637, 638, 7, 7,
	2, 2, 638, 639, 7, 13, 2, 2, 639, 640, 5, 18, 10, 2, 640, 641, 7, 14, 2,
	2, 641, 643, 3, 2, 2, 2, 642, 614, 3, 2, 2, 2, 642, 624, 3, 2, 2, 2, 642,
	632, 3, 2, 2, 2, 643, 121, 3, 2, 2, 2, 644, 645, 9, 4, 2, 2, 645, 123,
	3, 2, 2, 2, 646, 650, 7, 61, 2, 2, 647, 648, 7, 60, 2, 2, 648, 650, 7,
	61, 2, 2, 649, 646, 3, 2, 2, 2, 649, 647, 3, 2, 2, 2, 650, 125, 3, 2, 2,
	2, 651, 652, 9, 5, 2, 2, 652, 127, 3, 2, 2, 2, 653, 654, 9, 6, 2, 2, 654,
	129, 3, 2, 2, 2, 655, 656, 7, 30, 2, 2, 656, 131, 3, 2, 2, 2, 657, 658,
	7, 31, 2, 2, 658, 133, 3, 2, 2, 2, 659, 660, 9, 7, 2, 2, 660, 135, 3, 2,
	2, 2, 661, 662, 9, 8, 2, 2, 662, 137, 3, 2, 2, 2, 663, 664, 9, 9, 2, 2,
	664, 139, 3, 2, 2, 2, 62, 143, 158, 166, 170, 174, 179, 187, 193, 200,
	208, 215, 220, 233, 239, 243, 247, 251, 260, 264, 272, 277, 297, 308, 317,
	330, 332, 345, 352, 357, 365, 382, 392, 398, 402, 412, 415, 418, 436, 441,
	453, 464, 476, 488, 495, 500, 508, 513, 521, 526, 529, 540, 543, 564, 578,
	605, 609, 611, 617, 642, 649,
}
var literalNames = []string{
	"", "", "", "", "", "':'", "';'", "'.'", "','", "'['", "']'", "'('", "')'",
	"'{'", "'}'", "'>'", "'<'", "'=='", "'>='", "'<='", "'!='", "'*'", "'/'",
	"'%'", "'+'", "'-'", "'--'", "'++'", "", "", "", "'='", "'?'", "'!~'",
	"'=~'", "'FOR'", "'RETURN'", "'WAITFOR'", "'DISTINCT'", "'FILTER'", "'SORT'",
	"'LIMIT'", "'LET'", "'COLLECT'", "", "'NONE'", "'NULL'", "", "'USE'", "'INTO'",
	"'KEEP'", "'WITH'", "'COUNT'", "'ALL'", "'ANY'", "'AGGREGATE'", "'EVENT'",
	"'LIKE'", "", "'IN'", "'DO'", "'WHILE'", "'@'",
}
var symbolicNames = []string{
	"", "MultiLineComment", "SingleLineComment", "WhiteSpaces", "LineTerminator",
	"Colon", "SemiColon", "Dot", "Comma", "OpenBracket", "CloseBracket", "OpenParen",
	"CloseParen", "OpenBrace", "CloseBrace", "Gt", "Lt", "Eq", "Gte", "Lte",
	"Neq", "Multi", "Div", "Mod", "Plus", "Minus", "MinusMinus", "PlusPlus",
	"And", "Or", "Range", "Assign", "QuestionMark", "RegexNotMatch", "RegexMatch",
	"For", "Return", "Waitfor", "Distinct", "Filter", "Sort", "Limit", "Let",
	"Collect", "SortDirection", "None", "Null", "BooleanLiteral", "Use", "Into",
	"Keep", "With", "Count", "All", "Any", "Aggregate", "Event", "Like", "Not",
	"In", "Do", "While", "Param", "Identifier", "StringLiteral", "IntegerLiteral",
	"FloatLiteral", "NamespaceSegment",
}

var ruleNames = []string{
	"program", "head", "useExpression", "use", "body", "bodyStatement", "bodyExpression",
	"returnExpression", "forExpression", "forExpressionValueVariable", "forExpressionKeyVariable",
	"forExpressionSource", "forExpressionClause", "forExpressionStatement",
	"forExpressionBody", "forExpressionReturn", "filterClause", "limitClause",
	"limitClauseValue", "sortClause", "sortClauseExpression", "collectClause",
	"collectSelector", "collectGrouping", "collectAggregator", "collectAggregateSelector",
	"collectGroupVariable", "collectCounter", "waitForStatement", "waitForTimeout",
	"waitForEventName", "waitForEventSource", "waitForEventStatement", "variableDeclaration",
	"param", "variable", "rangeOperator", "arrayLiteral", "objectLiteral",
	"booleanLiteral", "stringLiteral", "integerLiteral", "floatLiteral", "noneLiteral",
	"arrayElementList", "propertyAssignment", "shorthandPropertyName", "computedPropertyName",
	"propertyName", "expressionGroup", "namespaceIdentifier", "namespace",
	"functionIdentifier", "functionCallExpression", "member", "memberPath",
	"memberExpression", "arguments", "expression", "forTernaryExpression",
	"arrayOperator", "inOperator", "equalityOperator", "regexpOperator", "logicalAndOperator",
	"logicalOrOperator", "multiplicativeOperator", "additiveOperator", "unaryOperator",
}

type FqlParser struct {
	*antlr.BaseParser
}

// NewFqlParser produces a new parser instance for the optional input antlr.TokenStream.
//
// The *FqlParser instance produced may be reused by calling the SetInputStream method.
// The initial parser configuration is expensive to construct, and the object is not thread-safe;
// however, if used within a Golang sync.Pool, the construction cost amortizes well and the
// objects can be used in a thread-safe manner.
func NewFqlParser(input antlr.TokenStream) *FqlParser {
	this := new(FqlParser)
	deserializer := antlr.NewATNDeserializer(nil)
	deserializedATN := deserializer.DeserializeFromUInt16(parserATN)
	decisionToDFA := make([]*antlr.DFA, len(deserializedATN.DecisionToState))
	for index, ds := range deserializedATN.DecisionToState {
		decisionToDFA[index] = antlr.NewDFA(ds, index)
	}
	this.BaseParser = antlr.NewBaseParser(input)

	this.Interpreter = antlr.NewParserATNSimulator(this, deserializedATN, decisionToDFA, antlr.NewPredictionContextCache())
	this.RuleNames = ruleNames
	this.LiteralNames = literalNames
	this.SymbolicNames = symbolicNames
	this.GrammarFileName = "FqlParser.g4"

	return this
}

// FqlParser tokens.
const (
	FqlParserEOF               = antlr.TokenEOF
	FqlParserMultiLineComment  = 1
	FqlParserSingleLineComment = 2
	FqlParserWhiteSpaces       = 3
	FqlParserLineTerminator    = 4
	FqlParserColon             = 5
	FqlParserSemiColon         = 6
	FqlParserDot               = 7
	FqlParserComma             = 8
	FqlParserOpenBracket       = 9
	FqlParserCloseBracket      = 10
	FqlParserOpenParen         = 11
	FqlParserCloseParen        = 12
	FqlParserOpenBrace         = 13
	FqlParserCloseBrace        = 14
	FqlParserGt                = 15
	FqlParserLt                = 16
	FqlParserEq                = 17
	FqlParserGte               = 18
	FqlParserLte               = 19
	FqlParserNeq               = 20
	FqlParserMulti             = 21
	FqlParserDiv               = 22
	FqlParserMod               = 23
	FqlParserPlus              = 24
	FqlParserMinus             = 25
	FqlParserMinusMinus        = 26
	FqlParserPlusPlus          = 27
	FqlParserAnd               = 28
	FqlParserOr                = 29
	FqlParserRange             = 30
	FqlParserAssign            = 31
	FqlParserQuestionMark      = 32
	FqlParserRegexNotMatch     = 33
	FqlParserRegexMatch        = 34
	FqlParserFor               = 35
	FqlParserReturn            = 36
	FqlParserWaitfor           = 37
	FqlParserDistinct          = 38
	FqlParserFilter            = 39
	FqlParserSort              = 40
	FqlParserLimit             = 41
	FqlParserLet               = 42
	FqlParserCollect           = 43
	FqlParserSortDirection     = 44
	FqlParserNone              = 45
	FqlParserNull              = 46
	FqlParserBooleanLiteral    = 47
	FqlParserUse               = 48
	FqlParserInto              = 49
	FqlParserKeep              = 50
	FqlParserWith              = 51
	FqlParserCount             = 52
	FqlParserAll               = 53
	FqlParserAny               = 54
	FqlParserAggregate         = 55
	FqlParserEvent             = 56
	FqlParserLike              = 57
	FqlParserNot               = 58
	FqlParserIn                = 59
	FqlParserDo                = 60
	FqlParserWhile             = 61
	FqlParserParam             = 62
	FqlParserIdentifier        = 63
	FqlParserStringLiteral     = 64
	FqlParserIntegerLiteral    = 65
	FqlParserFloatLiteral      = 66
	FqlParserNamespaceSegment  = 67
)

// FqlParser rules.
const (
	FqlParserRULE_program                    = 0
	FqlParserRULE_head                       = 1
	FqlParserRULE_useExpression              = 2
	FqlParserRULE_use                        = 3
	FqlParserRULE_body                       = 4
	FqlParserRULE_bodyStatement              = 5
	FqlParserRULE_bodyExpression             = 6
	FqlParserRULE_returnExpression           = 7
	FqlParserRULE_forExpression              = 8
	FqlParserRULE_forExpressionValueVariable = 9
	FqlParserRULE_forExpressionKeyVariable   = 10
	FqlParserRULE_forExpressionSource        = 11
	FqlParserRULE_forExpressionClause        = 12
	FqlParserRULE_forExpressionStatement     = 13
	FqlParserRULE_forExpressionBody          = 14
	FqlParserRULE_forExpressionReturn        = 15
	FqlParserRULE_filterClause               = 16
	FqlParserRULE_limitClause                = 17
	FqlParserRULE_limitClauseValue           = 18
	FqlParserRULE_sortClause                 = 19
	FqlParserRULE_sortClauseExpression       = 20
	FqlParserRULE_collectClause              = 21
	FqlParserRULE_collectSelector            = 22
	FqlParserRULE_collectGrouping            = 23
	FqlParserRULE_collectAggregator          = 24
	FqlParserRULE_collectAggregateSelector   = 25
	FqlParserRULE_collectGroupVariable       = 26
	FqlParserRULE_collectCounter             = 27
	FqlParserRULE_waitForStatement           = 28
	FqlParserRULE_waitForTimeout             = 29
	FqlParserRULE_waitForEventName           = 30
	FqlParserRULE_waitForEventSource         = 31
	FqlParserRULE_waitForEventStatement      = 32
	FqlParserRULE_variableDeclaration        = 33
	FqlParserRULE_param                      = 34
	FqlParserRULE_variable                   = 35
	FqlParserRULE_rangeOperator              = 36
	FqlParserRULE_arrayLiteral               = 37
	FqlParserRULE_objectLiteral              = 38
	FqlParserRULE_booleanLiteral             = 39
	FqlParserRULE_stringLiteral              = 40
	FqlParserRULE_integerLiteral             = 41
	FqlParserRULE_floatLiteral               = 42
	FqlParserRULE_noneLiteral                = 43
	FqlParserRULE_arrayElementList           = 44
	FqlParserRULE_propertyAssignment         = 45
	FqlParserRULE_shorthandPropertyName      = 46
	FqlParserRULE_computedPropertyName       = 47
	FqlParserRULE_propertyName               = 48
	FqlParserRULE_expressionGroup            = 49
	FqlParserRULE_namespaceIdentifier        = 50
	FqlParserRULE_namespace                  = 51
	FqlParserRULE_functionIdentifier         = 52
	FqlParserRULE_functionCallExpression     = 53
	FqlParserRULE_member                     = 54
	FqlParserRULE_memberPath                 = 55
	FqlParserRULE_memberExpression           = 56
	FqlParserRULE_arguments                  = 57
	FqlParserRULE_expression                 = 58
	FqlParserRULE_forTernaryExpression       = 59
	FqlParserRULE_arrayOperator              = 60
	FqlParserRULE_inOperator                 = 61
	FqlParserRULE_equalityOperator           = 62
	FqlParserRULE_regexpOperator             = 63
	FqlParserRULE_logicalAndOperator         = 64
	FqlParserRULE_logicalOrOperator          = 65
	FqlParserRULE_multiplicativeOperator     = 66
	FqlParserRULE_additiveOperator           = 67
	FqlParserRULE_unaryOperator              = 68
)

// IProgramContext is an interface to support dynamic dispatch.
type IProgramContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsProgramContext differentiates from other interfaces.
	IsProgramContext()
}

type ProgramContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyProgramContext() *ProgramContext {
	var p = new(ProgramContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = FqlParserRULE_program
	return p
}

func (*ProgramContext) IsProgramContext() {}

func NewProgramContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ProgramContext {
	var p = new(ProgramContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_program

	return p
}

func (s *ProgramContext) GetParser() antlr.Parser { return s.parser }

func (s *ProgramContext) Body() IBodyContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IBodyContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IBodyContext)
}

func (s *ProgramContext) AllHead() []IHeadContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IHeadContext)(nil)).Elem())
	var tst = make([]IHeadContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IHeadContext)
		}
	}

	return tst
}

func (s *ProgramContext) Head(i int) IHeadContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IHeadContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IHeadContext)
}

func (s *ProgramContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ProgramContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ProgramContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterProgram(s)
	}
}

func (s *ProgramContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitProgram(s)
	}
}

func (s *ProgramContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitProgram(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) Program() (localctx IProgramContext) {
	localctx = NewProgramContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 0, FqlParserRULE_program)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	var _alt int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(141)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 0, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			{
				p.SetState(138)
				p.Head()
			}

		}
		p.SetState(143)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 0, p.GetParserRuleContext())
	}
	{
		p.SetState(144)
		p.Body()
	}

	return localctx
}

// IHeadContext is an interface to support dynamic dispatch.
type IHeadContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsHeadContext differentiates from other interfaces.
	IsHeadContext()
}

type HeadContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyHeadContext() *HeadContext {
	var p = new(HeadContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = FqlParserRULE_head
	return p
}

func (*HeadContext) IsHeadContext() {}

func NewHeadContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *HeadContext {
	var p = new(HeadContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_head

	return p
}

func (s *HeadContext) GetParser() antlr.Parser { return s.parser }

func (s *HeadContext) UseExpression() IUseExpressionContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IUseExpressionContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IUseExpressionContext)
}

func (s *HeadContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *HeadContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *HeadContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterHead(s)
	}
}

func (s *HeadContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitHead(s)
	}
}

func (s *HeadContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitHead(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) Head() (localctx IHeadContext) {
	localctx = NewHeadContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 2, FqlParserRULE_head)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(146)
		p.UseExpression()
	}

	return localctx
}

// IUseExpressionContext is an interface to support dynamic dispatch.
type IUseExpressionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsUseExpressionContext differentiates from other interfaces.
	IsUseExpressionContext()
}

type UseExpressionContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyUseExpressionContext() *UseExpressionContext {
	var p = new(UseExpressionContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = FqlParserRULE_useExpression
	return p
}

func (*UseExpressionContext) IsUseExpressionContext() {}

func NewUseExpressionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *UseExpressionContext {
	var p = new(UseExpressionContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_useExpression

	return p
}

func (s *UseExpressionContext) GetParser() antlr.Parser { return s.parser }

func (s *UseExpressionContext) Use() IUseContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IUseContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IUseContext)
}

func (s *UseExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *UseExpressionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *UseExpressionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterUseExpression(s)
	}
}

func (s *UseExpressionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitUseExpression(s)
	}
}

func (s *UseExpressionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitUseExpression(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) UseExpression() (localctx IUseExpressionContext) {
	localctx = NewUseExpressionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 4, FqlParserRULE_useExpression)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(148)
		p.Use()
	}

	return localctx
}

// IUseContext is an interface to support dynamic dispatch.
type IUseContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsUseContext differentiates from other interfaces.
	IsUseContext()
}

type UseContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyUseContext() *UseContext {
	var p = new(UseContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = FqlParserRULE_use
	return p
}

func (*UseContext) IsUseContext() {}

func NewUseContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *UseContext {
	var p = new(UseContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_use

	return p
}

func (s *UseContext) GetParser() antlr.Parser { return s.parser }

func (s *UseContext) Use() antlr.TerminalNode {
	return s.GetToken(FqlParserUse, 0)
}

func (s *UseContext) NamespaceIdentifier() INamespaceIdentifierContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*INamespaceIdentifierContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(INamespaceIdentifierContext)
}

func (s *UseContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *UseContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *UseContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterUse(s)
	}
}

func (s *UseContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitUse(s)
	}
}

func (s *UseContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitUse(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) Use() (localctx IUseContext) {
	localctx = NewUseContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 6, FqlParserRULE_use)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(150)
		p.Match(FqlParserUse)
	}
	{
		p.SetState(151)
		p.NamespaceIdentifier()
	}

	return localctx
}

// IBodyContext is an interface to support dynamic dispatch.
type IBodyContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsBodyContext differentiates from other interfaces.
	IsBodyContext()
}

type BodyContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyBodyContext() *BodyContext {
	var p = new(BodyContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = FqlParserRULE_body
	return p
}

func (*BodyContext) IsBodyContext() {}

func NewBodyContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *BodyContext {
	var p = new(BodyContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_body

	return p
}

func (s *BodyContext) GetParser() antlr.Parser { return s.parser }

func (s *BodyContext) BodyExpression() IBodyExpressionContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IBodyExpressionContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IBodyExpressionContext)
}

func (s *BodyContext) AllBodyStatement() []IBodyStatementContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IBodyStatementContext)(nil)).Elem())
	var tst = make([]IBodyStatementContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IBodyStatementContext)
		}
	}

	return tst
}

func (s *BodyContext) BodyStatement(i int) IBodyStatementContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IBodyStatementContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IBodyStatementContext)
}

func (s *BodyContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *BodyContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *BodyContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterBody(s)
	}
}

func (s *BodyContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitBody(s)
	}
}

func (s *BodyContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitBody(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) Body() (localctx IBodyContext) {
	localctx = NewBodyContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 8, FqlParserRULE_body)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	var _alt int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(156)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 1, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			{
				p.SetState(153)
				p.BodyStatement()
			}

		}
		p.SetState(158)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 1, p.GetParserRuleContext())
	}
	{
		p.SetState(159)
		p.BodyExpression()
	}

	return localctx
}

// IBodyStatementContext is an interface to support dynamic dispatch.
type IBodyStatementContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsBodyStatementContext differentiates from other interfaces.
	IsBodyStatementContext()
}

type BodyStatementContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyBodyStatementContext() *BodyStatementContext {
	var p = new(BodyStatementContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = FqlParserRULE_bodyStatement
	return p
}

func (*BodyStatementContext) IsBodyStatementContext() {}

func NewBodyStatementContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *BodyStatementContext {
	var p = new(BodyStatementContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_bodyStatement

	return p
}

func (s *BodyStatementContext) GetParser() antlr.Parser { return s.parser }

func (s *BodyStatementContext) FunctionCallExpression() IFunctionCallExpressionContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IFunctionCallExpressionContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IFunctionCallExpressionContext)
}

func (s *BodyStatementContext) VariableDeclaration() IVariableDeclarationContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IVariableDeclarationContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IVariableDeclarationContext)
}

func (s *BodyStatementContext) WaitForStatement() IWaitForStatementContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IWaitForStatementContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IWaitForStatementContext)
}

func (s *BodyStatementContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *BodyStatementContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *BodyStatementContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterBodyStatement(s)
	}
}

func (s *BodyStatementContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitBodyStatement(s)
	}
}

func (s *BodyStatementContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitBodyStatement(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) BodyStatement() (localctx IBodyStatementContext) {
	localctx = NewBodyStatementContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 10, FqlParserRULE_bodyStatement)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(164)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 2, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(161)
			p.FunctionCallExpression()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(162)
			p.VariableDeclaration()
		}

	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(163)
			p.WaitForStatement()
		}

	}

	return localctx
}

// IBodyExpressionContext is an interface to support dynamic dispatch.
type IBodyExpressionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsBodyExpressionContext differentiates from other interfaces.
	IsBodyExpressionContext()
}

type BodyExpressionContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyBodyExpressionContext() *BodyExpressionContext {
	var p = new(BodyExpressionContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = FqlParserRULE_bodyExpression
	return p
}

func (*BodyExpressionContext) IsBodyExpressionContext() {}

func NewBodyExpressionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *BodyExpressionContext {
	var p = new(BodyExpressionContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_bodyExpression

	return p
}

func (s *BodyExpressionContext) GetParser() antlr.Parser { return s.parser }

func (s *BodyExpressionContext) ReturnExpression() IReturnExpressionContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IReturnExpressionContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IReturnExpressionContext)
}

func (s *BodyExpressionContext) ForExpression() IForExpressionContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IForExpressionContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IForExpressionContext)
}

func (s *BodyExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *BodyExpressionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *BodyExpressionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterBodyExpression(s)
	}
}

func (s *BodyExpressionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitBodyExpression(s)
	}
}

func (s *BodyExpressionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitBodyExpression(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) BodyExpression() (localctx IBodyExpressionContext) {
	localctx = NewBodyExpressionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 12, FqlParserRULE_bodyExpression)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(168)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case FqlParserReturn:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(166)
			p.ReturnExpression()
		}

	case FqlParserFor:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(167)
			p.ForExpression()
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}

	return localctx
}

// IReturnExpressionContext is an interface to support dynamic dispatch.
type IReturnExpressionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsReturnExpressionContext differentiates from other interfaces.
	IsReturnExpressionContext()
}

type ReturnExpressionContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyReturnExpressionContext() *ReturnExpressionContext {
	var p = new(ReturnExpressionContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = FqlParserRULE_returnExpression
	return p
}

func (*ReturnExpressionContext) IsReturnExpressionContext() {}

func NewReturnExpressionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ReturnExpressionContext {
	var p = new(ReturnExpressionContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_returnExpression

	return p
}

func (s *ReturnExpressionContext) GetParser() antlr.Parser { return s.parser }

func (s *ReturnExpressionContext) Return() antlr.TerminalNode {
	return s.GetToken(FqlParserReturn, 0)
}

func (s *ReturnExpressionContext) Expression() IExpressionContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExpressionContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *ReturnExpressionContext) Distinct() antlr.TerminalNode {
	return s.GetToken(FqlParserDistinct, 0)
}

func (s *ReturnExpressionContext) OpenParen() antlr.TerminalNode {
	return s.GetToken(FqlParserOpenParen, 0)
}

func (s *ReturnExpressionContext) ForExpression() IForExpressionContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IForExpressionContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IForExpressionContext)
}

func (s *ReturnExpressionContext) CloseParen() antlr.TerminalNode {
	return s.GetToken(FqlParserCloseParen, 0)
}

func (s *ReturnExpressionContext) ForTernaryExpression() IForTernaryExpressionContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IForTernaryExpressionContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IForTernaryExpressionContext)
}

func (s *ReturnExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ReturnExpressionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ReturnExpressionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterReturnExpression(s)
	}
}

func (s *ReturnExpressionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitReturnExpression(s)
	}
}

func (s *ReturnExpressionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitReturnExpression(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) ReturnExpression() (localctx IReturnExpressionContext) {
	localctx = NewReturnExpressionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 14, FqlParserRULE_returnExpression)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(185)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 6, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(170)
			p.Match(FqlParserReturn)
		}
		p.SetState(172)
		p.GetErrorHandler().Sync(p)

		if p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 4, p.GetParserRuleContext()) == 1 {
			{
				p.SetState(171)
				p.Match(FqlParserDistinct)
			}

		}
		{
			p.SetState(174)
			p.expression(0)
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(175)
			p.Match(FqlParserReturn)
		}
		p.SetState(177)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		if _la == FqlParserDistinct {
			{
				p.SetState(176)
				p.Match(FqlParserDistinct)
			}

		}
		{
			p.SetState(179)
			p.Match(FqlParserOpenParen)
		}
		{
			p.SetState(180)
			p.ForExpression()
		}
		{
			p.SetState(181)
			p.Match(FqlParserCloseParen)
		}

	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(183)
			p.Match(FqlParserReturn)
		}
		{
			p.SetState(184)
			p.ForTernaryExpression()
		}

	}

	return localctx
}

// IForExpressionContext is an interface to support dynamic dispatch.
type IForExpressionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsForExpressionContext differentiates from other interfaces.
	IsForExpressionContext()
}

type ForExpressionContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyForExpressionContext() *ForExpressionContext {
	var p = new(ForExpressionContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = FqlParserRULE_forExpression
	return p
}

func (*ForExpressionContext) IsForExpressionContext() {}

func NewForExpressionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ForExpressionContext {
	var p = new(ForExpressionContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_forExpression

	return p
}

func (s *ForExpressionContext) GetParser() antlr.Parser { return s.parser }

func (s *ForExpressionContext) For() antlr.TerminalNode {
	return s.GetToken(FqlParserFor, 0)
}

func (s *ForExpressionContext) ForExpressionValueVariable() IForExpressionValueVariableContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IForExpressionValueVariableContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IForExpressionValueVariableContext)
}

func (s *ForExpressionContext) In() antlr.TerminalNode {
	return s.GetToken(FqlParserIn, 0)
}

func (s *ForExpressionContext) ForExpressionSource() IForExpressionSourceContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IForExpressionSourceContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IForExpressionSourceContext)
}

func (s *ForExpressionContext) ForExpressionReturn() IForExpressionReturnContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IForExpressionReturnContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IForExpressionReturnContext)
}

func (s *ForExpressionContext) Comma() antlr.TerminalNode {
	return s.GetToken(FqlParserComma, 0)
}

func (s *ForExpressionContext) ForExpressionKeyVariable() IForExpressionKeyVariableContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IForExpressionKeyVariableContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IForExpressionKeyVariableContext)
}

func (s *ForExpressionContext) AllForExpressionBody() []IForExpressionBodyContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IForExpressionBodyContext)(nil)).Elem())
	var tst = make([]IForExpressionBodyContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IForExpressionBodyContext)
		}
	}

	return tst
}

func (s *ForExpressionContext) ForExpressionBody(i int) IForExpressionBodyContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IForExpressionBodyContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IForExpressionBodyContext)
}

func (s *ForExpressionContext) While() antlr.TerminalNode {
	return s.GetToken(FqlParserWhile, 0)
}

func (s *ForExpressionContext) Expression() IExpressionContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExpressionContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *ForExpressionContext) Do() antlr.TerminalNode {
	return s.GetToken(FqlParserDo, 0)
}

func (s *ForExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ForExpressionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ForExpressionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterForExpression(s)
	}
}

func (s *ForExpressionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitForExpression(s)
	}
}

func (s *ForExpressionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitForExpression(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) ForExpression() (localctx IForExpressionContext) {
	localctx = NewForExpressionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 16, FqlParserRULE_forExpression)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	var _alt int

	p.SetState(218)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 11, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(187)
			p.Match(FqlParserFor)
		}
		{
			p.SetState(188)
			p.ForExpressionValueVariable()
		}
		p.SetState(191)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		if _la == FqlParserComma {
			{
				p.SetState(189)
				p.Match(FqlParserComma)
			}
			{
				p.SetState(190)
				p.ForExpressionKeyVariable()
			}

		}
		{
			p.SetState(193)
			p.Match(FqlParserIn)
		}
		{
			p.SetState(194)
			p.ForExpressionSource()
		}
		p.SetState(198)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 8, p.GetParserRuleContext())

		for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
			if _alt == 1 {
				{
					p.SetState(195)
					p.ForExpressionBody()
				}

			}
			p.SetState(200)
			p.GetErrorHandler().Sync(p)
			_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 8, p.GetParserRuleContext())
		}
		{
			p.SetState(201)
			p.ForExpressionReturn()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(203)
			p.Match(FqlParserFor)
		}
		{
			p.SetState(204)
			p.ForExpressionValueVariable()
		}
		p.SetState(206)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		if _la == FqlParserDo {
			{
				p.SetState(205)
				p.Match(FqlParserDo)
			}

		}
		{
			p.SetState(208)
			p.Match(FqlParserWhile)
		}
		{
			p.SetState(209)
			p.expression(0)
		}
		p.SetState(213)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 10, p.GetParserRuleContext())

		for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
			if _alt == 1 {
				{
					p.SetState(210)
					p.ForExpressionBody()
				}

			}
			p.SetState(215)
			p.GetErrorHandler().Sync(p)
			_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 10, p.GetParserRuleContext())
		}
		{
			p.SetState(216)
			p.ForExpressionReturn()
		}

	}

	return localctx
}

// IForExpressionValueVariableContext is an interface to support dynamic dispatch.
type IForExpressionValueVariableContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsForExpressionValueVariableContext differentiates from other interfaces.
	IsForExpressionValueVariableContext()
}

type ForExpressionValueVariableContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyForExpressionValueVariableContext() *ForExpressionValueVariableContext {
	var p = new(ForExpressionValueVariableContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = FqlParserRULE_forExpressionValueVariable
	return p
}

func (*ForExpressionValueVariableContext) IsForExpressionValueVariableContext() {}

func NewForExpressionValueVariableContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ForExpressionValueVariableContext {
	var p = new(ForExpressionValueVariableContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_forExpressionValueVariable

	return p
}

func (s *ForExpressionValueVariableContext) GetParser() antlr.Parser { return s.parser }

func (s *ForExpressionValueVariableContext) Identifier() antlr.TerminalNode {
	return s.GetToken(FqlParserIdentifier, 0)
}

func (s *ForExpressionValueVariableContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ForExpressionValueVariableContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ForExpressionValueVariableContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterForExpressionValueVariable(s)
	}
}

func (s *ForExpressionValueVariableContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitForExpressionValueVariable(s)
	}
}

func (s *ForExpressionValueVariableContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitForExpressionValueVariable(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) ForExpressionValueVariable() (localctx IForExpressionValueVariableContext) {
	localctx = NewForExpressionValueVariableContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 18, FqlParserRULE_forExpressionValueVariable)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(220)
		p.Match(FqlParserIdentifier)
	}

	return localctx
}

// IForExpressionKeyVariableContext is an interface to support dynamic dispatch.
type IForExpressionKeyVariableContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsForExpressionKeyVariableContext differentiates from other interfaces.
	IsForExpressionKeyVariableContext()
}

type ForExpressionKeyVariableContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyForExpressionKeyVariableContext() *ForExpressionKeyVariableContext {
	var p = new(ForExpressionKeyVariableContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = FqlParserRULE_forExpressionKeyVariable
	return p
}

func (*ForExpressionKeyVariableContext) IsForExpressionKeyVariableContext() {}

func NewForExpressionKeyVariableContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ForExpressionKeyVariableContext {
	var p = new(ForExpressionKeyVariableContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_forExpressionKeyVariable

	return p
}

func (s *ForExpressionKeyVariableContext) GetParser() antlr.Parser { return s.parser }

func (s *ForExpressionKeyVariableContext) Identifier() antlr.TerminalNode {
	return s.GetToken(FqlParserIdentifier, 0)
}

func (s *ForExpressionKeyVariableContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ForExpressionKeyVariableContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ForExpressionKeyVariableContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterForExpressionKeyVariable(s)
	}
}

func (s *ForExpressionKeyVariableContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitForExpressionKeyVariable(s)
	}
}

func (s *ForExpressionKeyVariableContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitForExpressionKeyVariable(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) ForExpressionKeyVariable() (localctx IForExpressionKeyVariableContext) {
	localctx = NewForExpressionKeyVariableContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 20, FqlParserRULE_forExpressionKeyVariable)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(222)
		p.Match(FqlParserIdentifier)
	}

	return localctx
}

// IForExpressionSourceContext is an interface to support dynamic dispatch.
type IForExpressionSourceContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsForExpressionSourceContext differentiates from other interfaces.
	IsForExpressionSourceContext()
}

type ForExpressionSourceContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyForExpressionSourceContext() *ForExpressionSourceContext {
	var p = new(ForExpressionSourceContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = FqlParserRULE_forExpressionSource
	return p
}

func (*ForExpressionSourceContext) IsForExpressionSourceContext() {}

func NewForExpressionSourceContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ForExpressionSourceContext {
	var p = new(ForExpressionSourceContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_forExpressionSource

	return p
}

func (s *ForExpressionSourceContext) GetParser() antlr.Parser { return s.parser }

func (s *ForExpressionSourceContext) FunctionCallExpression() IFunctionCallExpressionContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IFunctionCallExpressionContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IFunctionCallExpressionContext)
}

func (s *ForExpressionSourceContext) ArrayLiteral() IArrayLiteralContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IArrayLiteralContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IArrayLiteralContext)
}

func (s *ForExpressionSourceContext) ObjectLiteral() IObjectLiteralContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IObjectLiteralContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IObjectLiteralContext)
}

func (s *ForExpressionSourceContext) Variable() IVariableContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IVariableContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IVariableContext)
}

func (s *ForExpressionSourceContext) MemberExpression() IMemberExpressionContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IMemberExpressionContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IMemberExpressionContext)
}

func (s *ForExpressionSourceContext) RangeOperator() IRangeOperatorContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IRangeOperatorContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IRangeOperatorContext)
}

func (s *ForExpressionSourceContext) Param() IParamContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IParamContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IParamContext)
}

func (s *ForExpressionSourceContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ForExpressionSourceContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ForExpressionSourceContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterForExpressionSource(s)
	}
}

func (s *ForExpressionSourceContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitForExpressionSource(s)
	}
}

func (s *ForExpressionSourceContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitForExpressionSource(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) ForExpressionSource() (localctx IForExpressionSourceContext) {
	localctx = NewForExpressionSourceContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 22, FqlParserRULE_forExpressionSource)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(231)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 12, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(224)
			p.FunctionCallExpression()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(225)
			p.ArrayLiteral()
		}

	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(226)
			p.ObjectLiteral()
		}

	case 4:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(227)
			p.Variable()
		}

	case 5:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(228)
			p.MemberExpression()
		}

	case 6:
		p.EnterOuterAlt(localctx, 6)
		{
			p.SetState(229)
			p.RangeOperator()
		}

	case 7:
		p.EnterOuterAlt(localctx, 7)
		{
			p.SetState(230)
			p.Param()
		}

	}

	return localctx
}

// IForExpressionClauseContext is an interface to support dynamic dispatch.
type IForExpressionClauseContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsForExpressionClauseContext differentiates from other interfaces.
	IsForExpressionClauseContext()
}

type ForExpressionClauseContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyForExpressionClauseContext() *ForExpressionClauseContext {
	var p = new(ForExpressionClauseContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = FqlParserRULE_forExpressionClause
	return p
}

func (*ForExpressionClauseContext) IsForExpressionClauseContext() {}

func NewForExpressionClauseContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ForExpressionClauseContext {
	var p = new(ForExpressionClauseContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_forExpressionClause

	return p
}

func (s *ForExpressionClauseContext) GetParser() antlr.Parser { return s.parser }

func (s *ForExpressionClauseContext) LimitClause() ILimitClauseContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ILimitClauseContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ILimitClauseContext)
}

func (s *ForExpressionClauseContext) SortClause() ISortClauseContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ISortClauseContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ISortClauseContext)
}

func (s *ForExpressionClauseContext) FilterClause() IFilterClauseContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IFilterClauseContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IFilterClauseContext)
}

func (s *ForExpressionClauseContext) CollectClause() ICollectClauseContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ICollectClauseContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ICollectClauseContext)
}

func (s *ForExpressionClauseContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ForExpressionClauseContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ForExpressionClauseContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterForExpressionClause(s)
	}
}

func (s *ForExpressionClauseContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitForExpressionClause(s)
	}
}

func (s *ForExpressionClauseContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitForExpressionClause(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) ForExpressionClause() (localctx IForExpressionClauseContext) {
	localctx = NewForExpressionClauseContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 24, FqlParserRULE_forExpressionClause)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(237)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case FqlParserLimit:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(233)
			p.LimitClause()
		}

	case FqlParserSort:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(234)
			p.SortClause()
		}

	case FqlParserFilter:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(235)
			p.FilterClause()
		}

	case FqlParserCollect:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(236)
			p.CollectClause()
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}

	return localctx
}

// IForExpressionStatementContext is an interface to support dynamic dispatch.
type IForExpressionStatementContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsForExpressionStatementContext differentiates from other interfaces.
	IsForExpressionStatementContext()
}

type ForExpressionStatementContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyForExpressionStatementContext() *ForExpressionStatementContext {
	var p = new(ForExpressionStatementContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = FqlParserRULE_forExpressionStatement
	return p
}

func (*ForExpressionStatementContext) IsForExpressionStatementContext() {}

func NewForExpressionStatementContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ForExpressionStatementContext {
	var p = new(ForExpressionStatementContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_forExpressionStatement

	return p
}

func (s *ForExpressionStatementContext) GetParser() antlr.Parser { return s.parser }

func (s *ForExpressionStatementContext) VariableDeclaration() IVariableDeclarationContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IVariableDeclarationContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IVariableDeclarationContext)
}

func (s *ForExpressionStatementContext) FunctionCallExpression() IFunctionCallExpressionContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IFunctionCallExpressionContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IFunctionCallExpressionContext)
}

func (s *ForExpressionStatementContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ForExpressionStatementContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ForExpressionStatementContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterForExpressionStatement(s)
	}
}

func (s *ForExpressionStatementContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitForExpressionStatement(s)
	}
}

func (s *ForExpressionStatementContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitForExpressionStatement(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) ForExpressionStatement() (localctx IForExpressionStatementContext) {
	localctx = NewForExpressionStatementContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 26, FqlParserRULE_forExpressionStatement)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(241)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 14, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(239)
			p.VariableDeclaration()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(240)
			p.FunctionCallExpression()
		}

	}

	return localctx
}

// IForExpressionBodyContext is an interface to support dynamic dispatch.
type IForExpressionBodyContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsForExpressionBodyContext differentiates from other interfaces.
	IsForExpressionBodyContext()
}

type ForExpressionBodyContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyForExpressionBodyContext() *ForExpressionBodyContext {
	var p = new(ForExpressionBodyContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = FqlParserRULE_forExpressionBody
	return p
}

func (*ForExpressionBodyContext) IsForExpressionBodyContext() {}

func NewForExpressionBodyContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ForExpressionBodyContext {
	var p = new(ForExpressionBodyContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_forExpressionBody

	return p
}

func (s *ForExpressionBodyContext) GetParser() antlr.Parser { return s.parser }

func (s *ForExpressionBodyContext) ForExpressionStatement() IForExpressionStatementContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IForExpressionStatementContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IForExpressionStatementContext)
}

func (s *ForExpressionBodyContext) ForExpressionClause() IForExpressionClauseContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IForExpressionClauseContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IForExpressionClauseContext)
}

func (s *ForExpressionBodyContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ForExpressionBodyContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ForExpressionBodyContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterForExpressionBody(s)
	}
}

func (s *ForExpressionBodyContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitForExpressionBody(s)
	}
}

func (s *ForExpressionBodyContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitForExpressionBody(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) ForExpressionBody() (localctx IForExpressionBodyContext) {
	localctx = NewForExpressionBodyContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 28, FqlParserRULE_forExpressionBody)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(245)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 15, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(243)
			p.ForExpressionStatement()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(244)
			p.ForExpressionClause()
		}

	}

	return localctx
}

// IForExpressionReturnContext is an interface to support dynamic dispatch.
type IForExpressionReturnContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsForExpressionReturnContext differentiates from other interfaces.
	IsForExpressionReturnContext()
}

type ForExpressionReturnContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyForExpressionReturnContext() *ForExpressionReturnContext {
	var p = new(ForExpressionReturnContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = FqlParserRULE_forExpressionReturn
	return p
}

func (*ForExpressionReturnContext) IsForExpressionReturnContext() {}

func NewForExpressionReturnContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ForExpressionReturnContext {
	var p = new(ForExpressionReturnContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_forExpressionReturn

	return p
}

func (s *ForExpressionReturnContext) GetParser() antlr.Parser { return s.parser }

func (s *ForExpressionReturnContext) ReturnExpression() IReturnExpressionContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IReturnExpressionContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IReturnExpressionContext)
}

func (s *ForExpressionReturnContext) ForExpression() IForExpressionContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IForExpressionContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IForExpressionContext)
}

func (s *ForExpressionReturnContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ForExpressionReturnContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ForExpressionReturnContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterForExpressionReturn(s)
	}
}

func (s *ForExpressionReturnContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitForExpressionReturn(s)
	}
}

func (s *ForExpressionReturnContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitForExpressionReturn(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) ForExpressionReturn() (localctx IForExpressionReturnContext) {
	localctx = NewForExpressionReturnContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 30, FqlParserRULE_forExpressionReturn)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(249)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case FqlParserReturn:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(247)
			p.ReturnExpression()
		}

	case FqlParserFor:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(248)
			p.ForExpression()
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}

	return localctx
}

// IFilterClauseContext is an interface to support dynamic dispatch.
type IFilterClauseContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsFilterClauseContext differentiates from other interfaces.
	IsFilterClauseContext()
}

type FilterClauseContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyFilterClauseContext() *FilterClauseContext {
	var p = new(FilterClauseContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = FqlParserRULE_filterClause
	return p
}

func (*FilterClauseContext) IsFilterClauseContext() {}

func NewFilterClauseContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FilterClauseContext {
	var p = new(FilterClauseContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_filterClause

	return p
}

func (s *FilterClauseContext) GetParser() antlr.Parser { return s.parser }

func (s *FilterClauseContext) Filter() antlr.TerminalNode {
	return s.GetToken(FqlParserFilter, 0)
}

func (s *FilterClauseContext) Expression() IExpressionContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExpressionContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *FilterClauseContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FilterClauseContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *FilterClauseContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterFilterClause(s)
	}
}

func (s *FilterClauseContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitFilterClause(s)
	}
}

func (s *FilterClauseContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitFilterClause(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) FilterClause() (localctx IFilterClauseContext) {
	localctx = NewFilterClauseContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 32, FqlParserRULE_filterClause)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(251)
		p.Match(FqlParserFilter)
	}
	{
		p.SetState(252)
		p.expression(0)
	}

	return localctx
}

// ILimitClauseContext is an interface to support dynamic dispatch.
type ILimitClauseContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsLimitClauseContext differentiates from other interfaces.
	IsLimitClauseContext()
}

type LimitClauseContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyLimitClauseContext() *LimitClauseContext {
	var p = new(LimitClauseContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = FqlParserRULE_limitClause
	return p
}

func (*LimitClauseContext) IsLimitClauseContext() {}

func NewLimitClauseContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *LimitClauseContext {
	var p = new(LimitClauseContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_limitClause

	return p
}

func (s *LimitClauseContext) GetParser() antlr.Parser { return s.parser }

func (s *LimitClauseContext) Limit() antlr.TerminalNode {
	return s.GetToken(FqlParserLimit, 0)
}

func (s *LimitClauseContext) AllLimitClauseValue() []ILimitClauseValueContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*ILimitClauseValueContext)(nil)).Elem())
	var tst = make([]ILimitClauseValueContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(ILimitClauseValueContext)
		}
	}

	return tst
}

func (s *LimitClauseContext) LimitClauseValue(i int) ILimitClauseValueContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ILimitClauseValueContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(ILimitClauseValueContext)
}

func (s *LimitClauseContext) Comma() antlr.TerminalNode {
	return s.GetToken(FqlParserComma, 0)
}

func (s *LimitClauseContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *LimitClauseContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *LimitClauseContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterLimitClause(s)
	}
}

func (s *LimitClauseContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitLimitClause(s)
	}
}

func (s *LimitClauseContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitLimitClause(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) LimitClause() (localctx ILimitClauseContext) {
	localctx = NewLimitClauseContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 34, FqlParserRULE_limitClause)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(254)
		p.Match(FqlParserLimit)
	}
	{
		p.SetState(255)
		p.LimitClauseValue()
	}
	p.SetState(258)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == FqlParserComma {
		{
			p.SetState(256)
			p.Match(FqlParserComma)
		}
		{
			p.SetState(257)
			p.LimitClauseValue()
		}

	}

	return localctx
}

// ILimitClauseValueContext is an interface to support dynamic dispatch.
type ILimitClauseValueContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsLimitClauseValueContext differentiates from other interfaces.
	IsLimitClauseValueContext()
}

type LimitClauseValueContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyLimitClauseValueContext() *LimitClauseValueContext {
	var p = new(LimitClauseValueContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = FqlParserRULE_limitClauseValue
	return p
}

func (*LimitClauseValueContext) IsLimitClauseValueContext() {}

func NewLimitClauseValueContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *LimitClauseValueContext {
	var p = new(LimitClauseValueContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_limitClauseValue

	return p
}

func (s *LimitClauseValueContext) GetParser() antlr.Parser { return s.parser }

func (s *LimitClauseValueContext) IntegerLiteral() antlr.TerminalNode {
	return s.GetToken(FqlParserIntegerLiteral, 0)
}

func (s *LimitClauseValueContext) Param() IParamContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IParamContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IParamContext)
}

func (s *LimitClauseValueContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *LimitClauseValueContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *LimitClauseValueContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterLimitClauseValue(s)
	}
}

func (s *LimitClauseValueContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitLimitClauseValue(s)
	}
}

func (s *LimitClauseValueContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitLimitClauseValue(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) LimitClauseValue() (localctx ILimitClauseValueContext) {
	localctx = NewLimitClauseValueContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 36, FqlParserRULE_limitClauseValue)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(262)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case FqlParserIntegerLiteral:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(260)
			p.Match(FqlParserIntegerLiteral)
		}

	case FqlParserParam:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(261)
			p.Param()
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}

	return localctx
}

// ISortClauseContext is an interface to support dynamic dispatch.
type ISortClauseContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsSortClauseContext differentiates from other interfaces.
	IsSortClauseContext()
}

type SortClauseContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptySortClauseContext() *SortClauseContext {
	var p = new(SortClauseContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = FqlParserRULE_sortClause
	return p
}

func (*SortClauseContext) IsSortClauseContext() {}

func NewSortClauseContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SortClauseContext {
	var p = new(SortClauseContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_sortClause

	return p
}

func (s *SortClauseContext) GetParser() antlr.Parser { return s.parser }

func (s *SortClauseContext) Sort() antlr.TerminalNode {
	return s.GetToken(FqlParserSort, 0)
}

func (s *SortClauseContext) AllSortClauseExpression() []ISortClauseExpressionContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*ISortClauseExpressionContext)(nil)).Elem())
	var tst = make([]ISortClauseExpressionContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(ISortClauseExpressionContext)
		}
	}

	return tst
}

func (s *SortClauseContext) SortClauseExpression(i int) ISortClauseExpressionContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ISortClauseExpressionContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(ISortClauseExpressionContext)
}

func (s *SortClauseContext) AllComma() []antlr.TerminalNode {
	return s.GetTokens(FqlParserComma)
}

func (s *SortClauseContext) Comma(i int) antlr.TerminalNode {
	return s.GetToken(FqlParserComma, i)
}

func (s *SortClauseContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SortClauseContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *SortClauseContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterSortClause(s)
	}
}

func (s *SortClauseContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitSortClause(s)
	}
}

func (s *SortClauseContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitSortClause(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) SortClause() (localctx ISortClauseContext) {
	localctx = NewSortClauseContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 38, FqlParserRULE_sortClause)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(264)
		p.Match(FqlParserSort)
	}
	{
		p.SetState(265)
		p.SortClauseExpression()
	}
	p.SetState(270)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == FqlParserComma {
		{
			p.SetState(266)
			p.Match(FqlParserComma)
		}
		{
			p.SetState(267)
			p.SortClauseExpression()
		}

		p.SetState(272)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}

	return localctx
}

// ISortClauseExpressionContext is an interface to support dynamic dispatch.
type ISortClauseExpressionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsSortClauseExpressionContext differentiates from other interfaces.
	IsSortClauseExpressionContext()
}

type SortClauseExpressionContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptySortClauseExpressionContext() *SortClauseExpressionContext {
	var p = new(SortClauseExpressionContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = FqlParserRULE_sortClauseExpression
	return p
}

func (*SortClauseExpressionContext) IsSortClauseExpressionContext() {}

func NewSortClauseExpressionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SortClauseExpressionContext {
	var p = new(SortClauseExpressionContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_sortClauseExpression

	return p
}

func (s *SortClauseExpressionContext) GetParser() antlr.Parser { return s.parser }

func (s *SortClauseExpressionContext) Expression() IExpressionContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExpressionContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *SortClauseExpressionContext) SortDirection() antlr.TerminalNode {
	return s.GetToken(FqlParserSortDirection, 0)
}

func (s *SortClauseExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SortClauseExpressionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *SortClauseExpressionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterSortClauseExpression(s)
	}
}

func (s *SortClauseExpressionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitSortClauseExpression(s)
	}
}

func (s *SortClauseExpressionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitSortClauseExpression(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) SortClauseExpression() (localctx ISortClauseExpressionContext) {
	localctx = NewSortClauseExpressionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 40, FqlParserRULE_sortClauseExpression)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(273)
		p.expression(0)
	}
	p.SetState(275)
	p.GetErrorHandler().Sync(p)

	if p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 20, p.GetParserRuleContext()) == 1 {
		{
			p.SetState(274)
			p.Match(FqlParserSortDirection)
		}

	}

	return localctx
}

// ICollectClauseContext is an interface to support dynamic dispatch.
type ICollectClauseContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsCollectClauseContext differentiates from other interfaces.
	IsCollectClauseContext()
}

type CollectClauseContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyCollectClauseContext() *CollectClauseContext {
	var p = new(CollectClauseContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = FqlParserRULE_collectClause
	return p
}

func (*CollectClauseContext) IsCollectClauseContext() {}

func NewCollectClauseContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *CollectClauseContext {
	var p = new(CollectClauseContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_collectClause

	return p
}

func (s *CollectClauseContext) GetParser() antlr.Parser { return s.parser }

func (s *CollectClauseContext) Collect() antlr.TerminalNode {
	return s.GetToken(FqlParserCollect, 0)
}

func (s *CollectClauseContext) CollectCounter() ICollectCounterContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ICollectCounterContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ICollectCounterContext)
}

func (s *CollectClauseContext) CollectAggregator() ICollectAggregatorContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ICollectAggregatorContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ICollectAggregatorContext)
}

func (s *CollectClauseContext) CollectGrouping() ICollectGroupingContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ICollectGroupingContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ICollectGroupingContext)
}

func (s *CollectClauseContext) CollectGroupVariable() ICollectGroupVariableContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ICollectGroupVariableContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ICollectGroupVariableContext)
}

func (s *CollectClauseContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *CollectClauseContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *CollectClauseContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterCollectClause(s)
	}
}

func (s *CollectClauseContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitCollectClause(s)
	}
}

func (s *CollectClauseContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitCollectClause(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) CollectClause() (localctx ICollectClauseContext) {
	localctx = NewCollectClauseContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 42, FqlParserRULE_collectClause)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(295)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 21, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(277)
			p.Match(FqlParserCollect)
		}
		{
			p.SetState(278)
			p.CollectCounter()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(279)
			p.Match(FqlParserCollect)
		}
		{
			p.SetState(280)
			p.CollectAggregator()
		}

	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(281)
			p.Match(FqlParserCollect)
		}
		{
			p.SetState(282)
			p.CollectGrouping()
		}
		{
			p.SetState(283)
			p.CollectAggregator()
		}

	case 4:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(285)
			p.Match(FqlParserCollect)
		}
		{
			p.SetState(286)
			p.CollectGrouping()
		}
		{
			p.SetState(287)
			p.CollectGroupVariable()
		}

	case 5:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(289)
			p.Match(FqlParserCollect)
		}
		{
			p.SetState(290)
			p.CollectGrouping()
		}
		{
			p.SetState(291)
			p.CollectCounter()
		}

	case 6:
		p.EnterOuterAlt(localctx, 6)
		{
			p.SetState(293)
			p.Match(FqlParserCollect)
		}
		{
			p.SetState(294)
			p.CollectGrouping()
		}

	}

	return localctx
}

// ICollectSelectorContext is an interface to support dynamic dispatch.
type ICollectSelectorContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsCollectSelectorContext differentiates from other interfaces.
	IsCollectSelectorContext()
}

type CollectSelectorContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyCollectSelectorContext() *CollectSelectorContext {
	var p = new(CollectSelectorContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = FqlParserRULE_collectSelector
	return p
}

func (*CollectSelectorContext) IsCollectSelectorContext() {}

func NewCollectSelectorContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *CollectSelectorContext {
	var p = new(CollectSelectorContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_collectSelector

	return p
}

func (s *CollectSelectorContext) GetParser() antlr.Parser { return s.parser }

func (s *CollectSelectorContext) Identifier() antlr.TerminalNode {
	return s.GetToken(FqlParserIdentifier, 0)
}

func (s *CollectSelectorContext) Assign() antlr.TerminalNode {
	return s.GetToken(FqlParserAssign, 0)
}

func (s *CollectSelectorContext) Expression() IExpressionContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExpressionContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *CollectSelectorContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *CollectSelectorContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *CollectSelectorContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterCollectSelector(s)
	}
}

func (s *CollectSelectorContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitCollectSelector(s)
	}
}

func (s *CollectSelectorContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitCollectSelector(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) CollectSelector() (localctx ICollectSelectorContext) {
	localctx = NewCollectSelectorContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 44, FqlParserRULE_collectSelector)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(297)
		p.Match(FqlParserIdentifier)
	}
	{
		p.SetState(298)
		p.Match(FqlParserAssign)
	}
	{
		p.SetState(299)
		p.expression(0)
	}

	return localctx
}

// ICollectGroupingContext is an interface to support dynamic dispatch.
type ICollectGroupingContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsCollectGroupingContext differentiates from other interfaces.
	IsCollectGroupingContext()
}

type CollectGroupingContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyCollectGroupingContext() *CollectGroupingContext {
	var p = new(CollectGroupingContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = FqlParserRULE_collectGrouping
	return p
}

func (*CollectGroupingContext) IsCollectGroupingContext() {}

func NewCollectGroupingContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *CollectGroupingContext {
	var p = new(CollectGroupingContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_collectGrouping

	return p
}

func (s *CollectGroupingContext) GetParser() antlr.Parser { return s.parser }

func (s *CollectGroupingContext) AllCollectSelector() []ICollectSelectorContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*ICollectSelectorContext)(nil)).Elem())
	var tst = make([]ICollectSelectorContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(ICollectSelectorContext)
		}
	}

	return tst
}

func (s *CollectGroupingContext) CollectSelector(i int) ICollectSelectorContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ICollectSelectorContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(ICollectSelectorContext)
}

func (s *CollectGroupingContext) AllComma() []antlr.TerminalNode {
	return s.GetTokens(FqlParserComma)
}

func (s *CollectGroupingContext) Comma(i int) antlr.TerminalNode {
	return s.GetToken(FqlParserComma, i)
}

func (s *CollectGroupingContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *CollectGroupingContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *CollectGroupingContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterCollectGrouping(s)
	}
}

func (s *CollectGroupingContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitCollectGrouping(s)
	}
}

func (s *CollectGroupingContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitCollectGrouping(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) CollectGrouping() (localctx ICollectGroupingContext) {
	localctx = NewCollectGroupingContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 46, FqlParserRULE_collectGrouping)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(301)
		p.CollectSelector()
	}
	p.SetState(306)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == FqlParserComma {
		{
			p.SetState(302)
			p.Match(FqlParserComma)
		}
		{
			p.SetState(303)
			p.CollectSelector()
		}

		p.SetState(308)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}

	return localctx
}

// ICollectAggregatorContext is an interface to support dynamic dispatch.
type ICollectAggregatorContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsCollectAggregatorContext differentiates from other interfaces.
	IsCollectAggregatorContext()
}

type CollectAggregatorContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyCollectAggregatorContext() *CollectAggregatorContext {
	var p = new(CollectAggregatorContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = FqlParserRULE_collectAggregator
	return p
}

func (*CollectAggregatorContext) IsCollectAggregatorContext() {}

func NewCollectAggregatorContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *CollectAggregatorContext {
	var p = new(CollectAggregatorContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_collectAggregator

	return p
}

func (s *CollectAggregatorContext) GetParser() antlr.Parser { return s.parser }

func (s *CollectAggregatorContext) Aggregate() antlr.TerminalNode {
	return s.GetToken(FqlParserAggregate, 0)
}

func (s *CollectAggregatorContext) AllCollectAggregateSelector() []ICollectAggregateSelectorContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*ICollectAggregateSelectorContext)(nil)).Elem())
	var tst = make([]ICollectAggregateSelectorContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(ICollectAggregateSelectorContext)
		}
	}

	return tst
}

func (s *CollectAggregatorContext) CollectAggregateSelector(i int) ICollectAggregateSelectorContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ICollectAggregateSelectorContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(ICollectAggregateSelectorContext)
}

func (s *CollectAggregatorContext) AllComma() []antlr.TerminalNode {
	return s.GetTokens(FqlParserComma)
}

func (s *CollectAggregatorContext) Comma(i int) antlr.TerminalNode {
	return s.GetToken(FqlParserComma, i)
}

func (s *CollectAggregatorContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *CollectAggregatorContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *CollectAggregatorContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterCollectAggregator(s)
	}
}

func (s *CollectAggregatorContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitCollectAggregator(s)
	}
}

func (s *CollectAggregatorContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitCollectAggregator(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) CollectAggregator() (localctx ICollectAggregatorContext) {
	localctx = NewCollectAggregatorContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 48, FqlParserRULE_collectAggregator)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(309)
		p.Match(FqlParserAggregate)
	}
	{
		p.SetState(310)
		p.CollectAggregateSelector()
	}
	p.SetState(315)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == FqlParserComma {
		{
			p.SetState(311)
			p.Match(FqlParserComma)
		}
		{
			p.SetState(312)
			p.CollectAggregateSelector()
		}

		p.SetState(317)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}

	return localctx
}

// ICollectAggregateSelectorContext is an interface to support dynamic dispatch.
type ICollectAggregateSelectorContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsCollectAggregateSelectorContext differentiates from other interfaces.
	IsCollectAggregateSelectorContext()
}

type CollectAggregateSelectorContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyCollectAggregateSelectorContext() *CollectAggregateSelectorContext {
	var p = new(CollectAggregateSelectorContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = FqlParserRULE_collectAggregateSelector
	return p
}

func (*CollectAggregateSelectorContext) IsCollectAggregateSelectorContext() {}

func NewCollectAggregateSelectorContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *CollectAggregateSelectorContext {
	var p = new(CollectAggregateSelectorContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_collectAggregateSelector

	return p
}

func (s *CollectAggregateSelectorContext) GetParser() antlr.Parser { return s.parser }

func (s *CollectAggregateSelectorContext) Identifier() antlr.TerminalNode {
	return s.GetToken(FqlParserIdentifier, 0)
}

func (s *CollectAggregateSelectorContext) Assign() antlr.TerminalNode {
	return s.GetToken(FqlParserAssign, 0)
}

func (s *CollectAggregateSelectorContext) FunctionCallExpression() IFunctionCallExpressionContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IFunctionCallExpressionContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IFunctionCallExpressionContext)
}

func (s *CollectAggregateSelectorContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *CollectAggregateSelectorContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *CollectAggregateSelectorContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterCollectAggregateSelector(s)
	}
}

func (s *CollectAggregateSelectorContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitCollectAggregateSelector(s)
	}
}

func (s *CollectAggregateSelectorContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitCollectAggregateSelector(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) CollectAggregateSelector() (localctx ICollectAggregateSelectorContext) {
	localctx = NewCollectAggregateSelectorContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 50, FqlParserRULE_collectAggregateSelector)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(318)
		p.Match(FqlParserIdentifier)
	}
	{
		p.SetState(319)
		p.Match(FqlParserAssign)
	}
	{
		p.SetState(320)
		p.FunctionCallExpression()
	}

	return localctx
}

// ICollectGroupVariableContext is an interface to support dynamic dispatch.
type ICollectGroupVariableContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsCollectGroupVariableContext differentiates from other interfaces.
	IsCollectGroupVariableContext()
}

type CollectGroupVariableContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyCollectGroupVariableContext() *CollectGroupVariableContext {
	var p = new(CollectGroupVariableContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = FqlParserRULE_collectGroupVariable
	return p
}

func (*CollectGroupVariableContext) IsCollectGroupVariableContext() {}

func NewCollectGroupVariableContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *CollectGroupVariableContext {
	var p = new(CollectGroupVariableContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_collectGroupVariable

	return p
}

func (s *CollectGroupVariableContext) GetParser() antlr.Parser { return s.parser }

func (s *CollectGroupVariableContext) Into() antlr.TerminalNode {
	return s.GetToken(FqlParserInto, 0)
}

func (s *CollectGroupVariableContext) CollectSelector() ICollectSelectorContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ICollectSelectorContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ICollectSelectorContext)
}

func (s *CollectGroupVariableContext) AllIdentifier() []antlr.TerminalNode {
	return s.GetTokens(FqlParserIdentifier)
}

func (s *CollectGroupVariableContext) Identifier(i int) antlr.TerminalNode {
	return s.GetToken(FqlParserIdentifier, i)
}

func (s *CollectGroupVariableContext) Keep() antlr.TerminalNode {
	return s.GetToken(FqlParserKeep, 0)
}

func (s *CollectGroupVariableContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *CollectGroupVariableContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *CollectGroupVariableContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterCollectGroupVariable(s)
	}
}

func (s *CollectGroupVariableContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitCollectGroupVariable(s)
	}
}

func (s *CollectGroupVariableContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitCollectGroupVariable(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) CollectGroupVariable() (localctx ICollectGroupVariableContext) {
	localctx = NewCollectGroupVariableContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 52, FqlParserRULE_collectGroupVariable)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(330)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 25, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(322)
			p.Match(FqlParserInto)
		}
		{
			p.SetState(323)
			p.CollectSelector()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(324)
			p.Match(FqlParserInto)
		}
		{
			p.SetState(325)
			p.Match(FqlParserIdentifier)
		}
		p.SetState(328)
		p.GetErrorHandler().Sync(p)

		if p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 24, p.GetParserRuleContext()) == 1 {
			{
				p.SetState(326)
				p.Match(FqlParserKeep)
			}
			{
				p.SetState(327)
				p.Match(FqlParserIdentifier)
			}

		}

	}

	return localctx
}

// ICollectCounterContext is an interface to support dynamic dispatch.
type ICollectCounterContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsCollectCounterContext differentiates from other interfaces.
	IsCollectCounterContext()
}

type CollectCounterContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyCollectCounterContext() *CollectCounterContext {
	var p = new(CollectCounterContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = FqlParserRULE_collectCounter
	return p
}

func (*CollectCounterContext) IsCollectCounterContext() {}

func NewCollectCounterContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *CollectCounterContext {
	var p = new(CollectCounterContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_collectCounter

	return p
}

func (s *CollectCounterContext) GetParser() antlr.Parser { return s.parser }

func (s *CollectCounterContext) With() antlr.TerminalNode {
	return s.GetToken(FqlParserWith, 0)
}

func (s *CollectCounterContext) Count() antlr.TerminalNode {
	return s.GetToken(FqlParserCount, 0)
}

func (s *CollectCounterContext) Into() antlr.TerminalNode {
	return s.GetToken(FqlParserInto, 0)
}

func (s *CollectCounterContext) Identifier() antlr.TerminalNode {
	return s.GetToken(FqlParserIdentifier, 0)
}

func (s *CollectCounterContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *CollectCounterContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *CollectCounterContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterCollectCounter(s)
	}
}

func (s *CollectCounterContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitCollectCounter(s)
	}
}

func (s *CollectCounterContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitCollectCounter(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) CollectCounter() (localctx ICollectCounterContext) {
	localctx = NewCollectCounterContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 54, FqlParserRULE_collectCounter)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(332)
		p.Match(FqlParserWith)
	}
	{
		p.SetState(333)
		p.Match(FqlParserCount)
	}
	{
		p.SetState(334)
		p.Match(FqlParserInto)
	}
	{
		p.SetState(335)
		p.Match(FqlParserIdentifier)
	}

	return localctx
}

// IWaitForStatementContext is an interface to support dynamic dispatch.
type IWaitForStatementContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsWaitForStatementContext differentiates from other interfaces.
	IsWaitForStatementContext()
}

type WaitForStatementContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyWaitForStatementContext() *WaitForStatementContext {
	var p = new(WaitForStatementContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = FqlParserRULE_waitForStatement
	return p
}

func (*WaitForStatementContext) IsWaitForStatementContext() {}

func NewWaitForStatementContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *WaitForStatementContext {
	var p = new(WaitForStatementContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_waitForStatement

	return p
}

func (s *WaitForStatementContext) GetParser() antlr.Parser { return s.parser }

func (s *WaitForStatementContext) WaitForEventStatement() IWaitForEventStatementContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IWaitForEventStatementContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IWaitForEventStatementContext)
}

func (s *WaitForStatementContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *WaitForStatementContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *WaitForStatementContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterWaitForStatement(s)
	}
}

func (s *WaitForStatementContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitWaitForStatement(s)
	}
}

func (s *WaitForStatementContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitWaitForStatement(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) WaitForStatement() (localctx IWaitForStatementContext) {
	localctx = NewWaitForStatementContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 56, FqlParserRULE_waitForStatement)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(337)
		p.WaitForEventStatement()
	}

	return localctx
}

// IWaitForTimeoutContext is an interface to support dynamic dispatch.
type IWaitForTimeoutContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsWaitForTimeoutContext differentiates from other interfaces.
	IsWaitForTimeoutContext()
}

type WaitForTimeoutContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyWaitForTimeoutContext() *WaitForTimeoutContext {
	var p = new(WaitForTimeoutContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = FqlParserRULE_waitForTimeout
	return p
}

func (*WaitForTimeoutContext) IsWaitForTimeoutContext() {}

func NewWaitForTimeoutContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *WaitForTimeoutContext {
	var p = new(WaitForTimeoutContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_waitForTimeout

	return p
}

func (s *WaitForTimeoutContext) GetParser() antlr.Parser { return s.parser }

func (s *WaitForTimeoutContext) IntegerLiteral() IIntegerLiteralContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IIntegerLiteralContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IIntegerLiteralContext)
}

func (s *WaitForTimeoutContext) Variable() IVariableContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IVariableContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IVariableContext)
}

func (s *WaitForTimeoutContext) FunctionCallExpression() IFunctionCallExpressionContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IFunctionCallExpressionContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IFunctionCallExpressionContext)
}

func (s *WaitForTimeoutContext) MemberExpression() IMemberExpressionContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IMemberExpressionContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IMemberExpressionContext)
}

func (s *WaitForTimeoutContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *WaitForTimeoutContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *WaitForTimeoutContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterWaitForTimeout(s)
	}
}

func (s *WaitForTimeoutContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitWaitForTimeout(s)
	}
}

func (s *WaitForTimeoutContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitWaitForTimeout(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) WaitForTimeout() (localctx IWaitForTimeoutContext) {
	localctx = NewWaitForTimeoutContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 58, FqlParserRULE_waitForTimeout)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(343)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 26, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(339)
			p.IntegerLiteral()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(340)
			p.Variable()
		}

	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(341)
			p.FunctionCallExpression()
		}

	case 4:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(342)
			p.MemberExpression()
		}

	}

	return localctx
}

// IWaitForEventNameContext is an interface to support dynamic dispatch.
type IWaitForEventNameContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsWaitForEventNameContext differentiates from other interfaces.
	IsWaitForEventNameContext()
}

type WaitForEventNameContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyWaitForEventNameContext() *WaitForEventNameContext {
	var p = new(WaitForEventNameContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = FqlParserRULE_waitForEventName
	return p
}

func (*WaitForEventNameContext) IsWaitForEventNameContext() {}

func NewWaitForEventNameContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *WaitForEventNameContext {
	var p = new(WaitForEventNameContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_waitForEventName

	return p
}

func (s *WaitForEventNameContext) GetParser() antlr.Parser { return s.parser }

func (s *WaitForEventNameContext) StringLiteral() IStringLiteralContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IStringLiteralContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IStringLiteralContext)
}

func (s *WaitForEventNameContext) Variable() IVariableContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IVariableContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IVariableContext)
}

func (s *WaitForEventNameContext) Param() IParamContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IParamContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IParamContext)
}

func (s *WaitForEventNameContext) FunctionCallExpression() IFunctionCallExpressionContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IFunctionCallExpressionContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IFunctionCallExpressionContext)
}

func (s *WaitForEventNameContext) MemberExpression() IMemberExpressionContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IMemberExpressionContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IMemberExpressionContext)
}

func (s *WaitForEventNameContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *WaitForEventNameContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *WaitForEventNameContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterWaitForEventName(s)
	}
}

func (s *WaitForEventNameContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitWaitForEventName(s)
	}
}

func (s *WaitForEventNameContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitWaitForEventName(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) WaitForEventName() (localctx IWaitForEventNameContext) {
	localctx = NewWaitForEventNameContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 60, FqlParserRULE_waitForEventName)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(350)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 27, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(345)
			p.StringLiteral()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(346)
			p.Variable()
		}

	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(347)
			p.Param()
		}

	case 4:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(348)
			p.FunctionCallExpression()
		}

	case 5:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(349)
			p.MemberExpression()
		}

	}

	return localctx
}

// IWaitForEventSourceContext is an interface to support dynamic dispatch.
type IWaitForEventSourceContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsWaitForEventSourceContext differentiates from other interfaces.
	IsWaitForEventSourceContext()
}

type WaitForEventSourceContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyWaitForEventSourceContext() *WaitForEventSourceContext {
	var p = new(WaitForEventSourceContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = FqlParserRULE_waitForEventSource
	return p
}

func (*WaitForEventSourceContext) IsWaitForEventSourceContext() {}

func NewWaitForEventSourceContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *WaitForEventSourceContext {
	var p = new(WaitForEventSourceContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_waitForEventSource

	return p
}

func (s *WaitForEventSourceContext) GetParser() antlr.Parser { return s.parser }

func (s *WaitForEventSourceContext) Variable() IVariableContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IVariableContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IVariableContext)
}

func (s *WaitForEventSourceContext) FunctionCallExpression() IFunctionCallExpressionContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IFunctionCallExpressionContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IFunctionCallExpressionContext)
}

func (s *WaitForEventSourceContext) MemberExpression() IMemberExpressionContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IMemberExpressionContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IMemberExpressionContext)
}

func (s *WaitForEventSourceContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *WaitForEventSourceContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *WaitForEventSourceContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterWaitForEventSource(s)
	}
}

func (s *WaitForEventSourceContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitWaitForEventSource(s)
	}
}

func (s *WaitForEventSourceContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitWaitForEventSource(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) WaitForEventSource() (localctx IWaitForEventSourceContext) {
	localctx = NewWaitForEventSourceContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 62, FqlParserRULE_waitForEventSource)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(355)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 28, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(352)
			p.Variable()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(353)
			p.FunctionCallExpression()
		}

	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(354)
			p.MemberExpression()
		}

	}

	return localctx
}

// IWaitForEventStatementContext is an interface to support dynamic dispatch.
type IWaitForEventStatementContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsWaitForEventStatementContext differentiates from other interfaces.
	IsWaitForEventStatementContext()
}

type WaitForEventStatementContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyWaitForEventStatementContext() *WaitForEventStatementContext {
	var p = new(WaitForEventStatementContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = FqlParserRULE_waitForEventStatement
	return p
}

func (*WaitForEventStatementContext) IsWaitForEventStatementContext() {}

func NewWaitForEventStatementContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *WaitForEventStatementContext {
	var p = new(WaitForEventStatementContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_waitForEventStatement

	return p
}

func (s *WaitForEventStatementContext) GetParser() antlr.Parser { return s.parser }

func (s *WaitForEventStatementContext) Waitfor() antlr.TerminalNode {
	return s.GetToken(FqlParserWaitfor, 0)
}

func (s *WaitForEventStatementContext) Event() antlr.TerminalNode {
	return s.GetToken(FqlParserEvent, 0)
}

func (s *WaitForEventStatementContext) WaitForEventName() IWaitForEventNameContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IWaitForEventNameContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IWaitForEventNameContext)
}

func (s *WaitForEventStatementContext) In() antlr.TerminalNode {
	return s.GetToken(FqlParserIn, 0)
}

func (s *WaitForEventStatementContext) WaitForEventSource() IWaitForEventSourceContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IWaitForEventSourceContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IWaitForEventSourceContext)
}

func (s *WaitForEventStatementContext) WaitForTimeout() IWaitForTimeoutContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IWaitForTimeoutContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IWaitForTimeoutContext)
}

func (s *WaitForEventStatementContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *WaitForEventStatementContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *WaitForEventStatementContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterWaitForEventStatement(s)
	}
}

func (s *WaitForEventStatementContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitWaitForEventStatement(s)
	}
}

func (s *WaitForEventStatementContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitWaitForEventStatement(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) WaitForEventStatement() (localctx IWaitForEventStatementContext) {
	localctx = NewWaitForEventStatementContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 64, FqlParserRULE_waitForEventStatement)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(357)
		p.Match(FqlParserWaitfor)
	}
	{
		p.SetState(358)
		p.Match(FqlParserEvent)
	}
	{
		p.SetState(359)
		p.WaitForEventName()
	}
	{
		p.SetState(360)
		p.Match(FqlParserIn)
	}
	{
		p.SetState(361)
		p.WaitForEventSource()
	}
	p.SetState(363)
	p.GetErrorHandler().Sync(p)

	if p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 29, p.GetParserRuleContext()) == 1 {
		{
			p.SetState(362)
			p.WaitForTimeout()
		}

	}

	return localctx
}

// IVariableDeclarationContext is an interface to support dynamic dispatch.
type IVariableDeclarationContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsVariableDeclarationContext differentiates from other interfaces.
	IsVariableDeclarationContext()
}

type VariableDeclarationContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyVariableDeclarationContext() *VariableDeclarationContext {
	var p = new(VariableDeclarationContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = FqlParserRULE_variableDeclaration
	return p
}

func (*VariableDeclarationContext) IsVariableDeclarationContext() {}

func NewVariableDeclarationContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *VariableDeclarationContext {
	var p = new(VariableDeclarationContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_variableDeclaration

	return p
}

func (s *VariableDeclarationContext) GetParser() antlr.Parser { return s.parser }

func (s *VariableDeclarationContext) Let() antlr.TerminalNode {
	return s.GetToken(FqlParserLet, 0)
}

func (s *VariableDeclarationContext) Identifier() antlr.TerminalNode {
	return s.GetToken(FqlParserIdentifier, 0)
}

func (s *VariableDeclarationContext) Assign() antlr.TerminalNode {
	return s.GetToken(FqlParserAssign, 0)
}

func (s *VariableDeclarationContext) Expression() IExpressionContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExpressionContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *VariableDeclarationContext) OpenParen() antlr.TerminalNode {
	return s.GetToken(FqlParserOpenParen, 0)
}

func (s *VariableDeclarationContext) ForExpression() IForExpressionContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IForExpressionContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IForExpressionContext)
}

func (s *VariableDeclarationContext) CloseParen() antlr.TerminalNode {
	return s.GetToken(FqlParserCloseParen, 0)
}

func (s *VariableDeclarationContext) ForTernaryExpression() IForTernaryExpressionContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IForTernaryExpressionContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IForTernaryExpressionContext)
}

func (s *VariableDeclarationContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *VariableDeclarationContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *VariableDeclarationContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterVariableDeclaration(s)
	}
}

func (s *VariableDeclarationContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitVariableDeclaration(s)
	}
}

func (s *VariableDeclarationContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitVariableDeclaration(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) VariableDeclaration() (localctx IVariableDeclarationContext) {
	localctx = NewVariableDeclarationContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 66, FqlParserRULE_variableDeclaration)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(380)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 30, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(365)
			p.Match(FqlParserLet)
		}
		{
			p.SetState(366)
			p.Match(FqlParserIdentifier)
		}
		{
			p.SetState(367)
			p.Match(FqlParserAssign)
		}
		{
			p.SetState(368)
			p.expression(0)
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(369)
			p.Match(FqlParserLet)
		}
		{
			p.SetState(370)
			p.Match(FqlParserIdentifier)
		}
		{
			p.SetState(371)
			p.Match(FqlParserAssign)
		}
		{
			p.SetState(372)
			p.Match(FqlParserOpenParen)
		}
		{
			p.SetState(373)
			p.ForExpression()
		}
		{
			p.SetState(374)
			p.Match(FqlParserCloseParen)
		}

	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(376)
			p.Match(FqlParserLet)
		}
		{
			p.SetState(377)
			p.Match(FqlParserIdentifier)
		}
		{
			p.SetState(378)
			p.Match(FqlParserAssign)
		}
		{
			p.SetState(379)
			p.ForTernaryExpression()
		}

	}

	return localctx
}

// IParamContext is an interface to support dynamic dispatch.
type IParamContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsParamContext differentiates from other interfaces.
	IsParamContext()
}

type ParamContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyParamContext() *ParamContext {
	var p = new(ParamContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = FqlParserRULE_param
	return p
}

func (*ParamContext) IsParamContext() {}

func NewParamContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ParamContext {
	var p = new(ParamContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_param

	return p
}

func (s *ParamContext) GetParser() antlr.Parser { return s.parser }

func (s *ParamContext) Param() antlr.TerminalNode {
	return s.GetToken(FqlParserParam, 0)
}

func (s *ParamContext) Identifier() antlr.TerminalNode {
	return s.GetToken(FqlParserIdentifier, 0)
}

func (s *ParamContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ParamContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ParamContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterParam(s)
	}
}

func (s *ParamContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitParam(s)
	}
}

func (s *ParamContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitParam(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) Param() (localctx IParamContext) {
	localctx = NewParamContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 68, FqlParserRULE_param)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(382)
		p.Match(FqlParserParam)
	}
	{
		p.SetState(383)
		p.Match(FqlParserIdentifier)
	}

	return localctx
}

// IVariableContext is an interface to support dynamic dispatch.
type IVariableContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsVariableContext differentiates from other interfaces.
	IsVariableContext()
}

type VariableContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyVariableContext() *VariableContext {
	var p = new(VariableContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = FqlParserRULE_variable
	return p
}

func (*VariableContext) IsVariableContext() {}

func NewVariableContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *VariableContext {
	var p = new(VariableContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_variable

	return p
}

func (s *VariableContext) GetParser() antlr.Parser { return s.parser }

func (s *VariableContext) Identifier() antlr.TerminalNode {
	return s.GetToken(FqlParserIdentifier, 0)
}

func (s *VariableContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *VariableContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *VariableContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterVariable(s)
	}
}

func (s *VariableContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitVariable(s)
	}
}

func (s *VariableContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitVariable(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) Variable() (localctx IVariableContext) {
	localctx = NewVariableContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 70, FqlParserRULE_variable)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(385)
		p.Match(FqlParserIdentifier)
	}

	return localctx
}

// IRangeOperatorContext is an interface to support dynamic dispatch.
type IRangeOperatorContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsRangeOperatorContext differentiates from other interfaces.
	IsRangeOperatorContext()
}

type RangeOperatorContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyRangeOperatorContext() *RangeOperatorContext {
	var p = new(RangeOperatorContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = FqlParserRULE_rangeOperator
	return p
}

func (*RangeOperatorContext) IsRangeOperatorContext() {}

func NewRangeOperatorContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *RangeOperatorContext {
	var p = new(RangeOperatorContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_rangeOperator

	return p
}

func (s *RangeOperatorContext) GetParser() antlr.Parser { return s.parser }

func (s *RangeOperatorContext) Range() antlr.TerminalNode {
	return s.GetToken(FqlParserRange, 0)
}

func (s *RangeOperatorContext) AllIntegerLiteral() []IIntegerLiteralContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IIntegerLiteralContext)(nil)).Elem())
	var tst = make([]IIntegerLiteralContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IIntegerLiteralContext)
		}
	}

	return tst
}

func (s *RangeOperatorContext) IntegerLiteral(i int) IIntegerLiteralContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IIntegerLiteralContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IIntegerLiteralContext)
}

func (s *RangeOperatorContext) AllVariable() []IVariableContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IVariableContext)(nil)).Elem())
	var tst = make([]IVariableContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IVariableContext)
		}
	}

	return tst
}

func (s *RangeOperatorContext) Variable(i int) IVariableContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IVariableContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IVariableContext)
}

func (s *RangeOperatorContext) AllParam() []IParamContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IParamContext)(nil)).Elem())
	var tst = make([]IParamContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IParamContext)
		}
	}

	return tst
}

func (s *RangeOperatorContext) Param(i int) IParamContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IParamContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IParamContext)
}

func (s *RangeOperatorContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *RangeOperatorContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *RangeOperatorContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterRangeOperator(s)
	}
}

func (s *RangeOperatorContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitRangeOperator(s)
	}
}

func (s *RangeOperatorContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitRangeOperator(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) RangeOperator() (localctx IRangeOperatorContext) {
	localctx = NewRangeOperatorContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 72, FqlParserRULE_rangeOperator)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	p.SetState(390)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case FqlParserIntegerLiteral:
		{
			p.SetState(387)
			p.IntegerLiteral()
		}

	case FqlParserIdentifier:
		{
			p.SetState(388)
			p.Variable()
		}

	case FqlParserParam:
		{
			p.SetState(389)
			p.Param()
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}
	{
		p.SetState(392)
		p.Match(FqlParserRange)
	}
	p.SetState(396)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case FqlParserIntegerLiteral:
		{
			p.SetState(393)
			p.IntegerLiteral()
		}

	case FqlParserIdentifier:
		{
			p.SetState(394)
			p.Variable()
		}

	case FqlParserParam:
		{
			p.SetState(395)
			p.Param()
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}

	return localctx
}

// IArrayLiteralContext is an interface to support dynamic dispatch.
type IArrayLiteralContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsArrayLiteralContext differentiates from other interfaces.
	IsArrayLiteralContext()
}

type ArrayLiteralContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyArrayLiteralContext() *ArrayLiteralContext {
	var p = new(ArrayLiteralContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = FqlParserRULE_arrayLiteral
	return p
}

func (*ArrayLiteralContext) IsArrayLiteralContext() {}

func NewArrayLiteralContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ArrayLiteralContext {
	var p = new(ArrayLiteralContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_arrayLiteral

	return p
}

func (s *ArrayLiteralContext) GetParser() antlr.Parser { return s.parser }

func (s *ArrayLiteralContext) OpenBracket() antlr.TerminalNode {
	return s.GetToken(FqlParserOpenBracket, 0)
}

func (s *ArrayLiteralContext) CloseBracket() antlr.TerminalNode {
	return s.GetToken(FqlParserCloseBracket, 0)
}

func (s *ArrayLiteralContext) ArrayElementList() IArrayElementListContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IArrayElementListContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IArrayElementListContext)
}

func (s *ArrayLiteralContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ArrayLiteralContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ArrayLiteralContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterArrayLiteral(s)
	}
}

func (s *ArrayLiteralContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitArrayLiteral(s)
	}
}

func (s *ArrayLiteralContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitArrayLiteral(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) ArrayLiteral() (localctx IArrayLiteralContext) {
	localctx = NewArrayLiteralContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 74, FqlParserRULE_arrayLiteral)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(398)
		p.Match(FqlParserOpenBracket)
	}
	p.SetState(400)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if (((_la-9)&-(0x1f+1)) == 0 && ((1<<uint((_la-9)))&((1<<(FqlParserOpenBracket-9))|(1<<(FqlParserOpenParen-9))|(1<<(FqlParserOpenBrace-9))|(1<<(FqlParserPlus-9))|(1<<(FqlParserMinus-9))|(1<<(FqlParserAnd-9))|(1<<(FqlParserOr-9))|(1<<(FqlParserFor-9))|(1<<(FqlParserReturn-9))|(1<<(FqlParserWaitfor-9))|(1<<(FqlParserDistinct-9))|(1<<(FqlParserFilter-9))|(1<<(FqlParserSort-9)))) != 0) || (((_la-41)&-(0x1f+1)) == 0 && ((1<<uint((_la-41)))&((1<<(FqlParserLimit-41))|(1<<(FqlParserLet-41))|(1<<(FqlParserCollect-41))|(1<<(FqlParserSortDirection-41))|(1<<(FqlParserNone-41))|(1<<(FqlParserNull-41))|(1<<(FqlParserBooleanLiteral-41))|(1<<(FqlParserUse-41))|(1<<(FqlParserInto-41))|(1<<(FqlParserKeep-41))|(1<<(FqlParserWith-41))|(1<<(FqlParserCount-41))|(1<<(FqlParserAll-41))|(1<<(FqlParserAny-41))|(1<<(FqlParserAggregate-41))|(1<<(FqlParserEvent-41))|(1<<(FqlParserLike-41))|(1<<(FqlParserNot-41))|(1<<(FqlParserIn-41))|(1<<(FqlParserParam-41))|(1<<(FqlParserIdentifier-41))|(1<<(FqlParserStringLiteral-41))|(1<<(FqlParserIntegerLiteral-41))|(1<<(FqlParserFloatLiteral-41))|(1<<(FqlParserNamespaceSegment-41)))) != 0) {
		{
			p.SetState(399)
			p.ArrayElementList()
		}

	}
	{
		p.SetState(402)
		p.Match(FqlParserCloseBracket)
	}

	return localctx
}

// IObjectLiteralContext is an interface to support dynamic dispatch.
type IObjectLiteralContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsObjectLiteralContext differentiates from other interfaces.
	IsObjectLiteralContext()
}

type ObjectLiteralContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyObjectLiteralContext() *ObjectLiteralContext {
	var p = new(ObjectLiteralContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = FqlParserRULE_objectLiteral
	return p
}

func (*ObjectLiteralContext) IsObjectLiteralContext() {}

func NewObjectLiteralContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ObjectLiteralContext {
	var p = new(ObjectLiteralContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_objectLiteral

	return p
}

func (s *ObjectLiteralContext) GetParser() antlr.Parser { return s.parser }

func (s *ObjectLiteralContext) OpenBrace() antlr.TerminalNode {
	return s.GetToken(FqlParserOpenBrace, 0)
}

func (s *ObjectLiteralContext) CloseBrace() antlr.TerminalNode {
	return s.GetToken(FqlParserCloseBrace, 0)
}

func (s *ObjectLiteralContext) AllPropertyAssignment() []IPropertyAssignmentContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IPropertyAssignmentContext)(nil)).Elem())
	var tst = make([]IPropertyAssignmentContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IPropertyAssignmentContext)
		}
	}

	return tst
}

func (s *ObjectLiteralContext) PropertyAssignment(i int) IPropertyAssignmentContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IPropertyAssignmentContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IPropertyAssignmentContext)
}

func (s *ObjectLiteralContext) AllComma() []antlr.TerminalNode {
	return s.GetTokens(FqlParserComma)
}

func (s *ObjectLiteralContext) Comma(i int) antlr.TerminalNode {
	return s.GetToken(FqlParserComma, i)
}

func (s *ObjectLiteralContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ObjectLiteralContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ObjectLiteralContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterObjectLiteral(s)
	}
}

func (s *ObjectLiteralContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitObjectLiteral(s)
	}
}

func (s *ObjectLiteralContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitObjectLiteral(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) ObjectLiteral() (localctx IObjectLiteralContext) {
	localctx = NewObjectLiteralContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 76, FqlParserRULE_objectLiteral)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	var _alt int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(404)
		p.Match(FqlParserOpenBrace)
	}
	p.SetState(413)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == FqlParserOpenBracket || (((_la-62)&-(0x1f+1)) == 0 && ((1<<uint((_la-62)))&((1<<(FqlParserParam-62))|(1<<(FqlParserIdentifier-62))|(1<<(FqlParserStringLiteral-62)))) != 0) {
		{
			p.SetState(405)
			p.PropertyAssignment()
		}
		p.SetState(410)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 34, p.GetParserRuleContext())

		for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
			if _alt == 1 {
				{
					p.SetState(406)
					p.Match(FqlParserComma)
				}
				{
					p.SetState(407)
					p.PropertyAssignment()
				}

			}
			p.SetState(412)
			p.GetErrorHandler().Sync(p)
			_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 34, p.GetParserRuleContext())
		}

	}
	p.SetState(416)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == FqlParserComma {
		{
			p.SetState(415)
			p.Match(FqlParserComma)
		}

	}
	{
		p.SetState(418)
		p.Match(FqlParserCloseBrace)
	}

	return localctx
}

// IBooleanLiteralContext is an interface to support dynamic dispatch.
type IBooleanLiteralContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsBooleanLiteralContext differentiates from other interfaces.
	IsBooleanLiteralContext()
}

type BooleanLiteralContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyBooleanLiteralContext() *BooleanLiteralContext {
	var p = new(BooleanLiteralContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = FqlParserRULE_booleanLiteral
	return p
}

func (*BooleanLiteralContext) IsBooleanLiteralContext() {}

func NewBooleanLiteralContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *BooleanLiteralContext {
	var p = new(BooleanLiteralContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_booleanLiteral

	return p
}

func (s *BooleanLiteralContext) GetParser() antlr.Parser { return s.parser }

func (s *BooleanLiteralContext) BooleanLiteral() antlr.TerminalNode {
	return s.GetToken(FqlParserBooleanLiteral, 0)
}

func (s *BooleanLiteralContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *BooleanLiteralContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *BooleanLiteralContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterBooleanLiteral(s)
	}
}

func (s *BooleanLiteralContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitBooleanLiteral(s)
	}
}

func (s *BooleanLiteralContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitBooleanLiteral(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) BooleanLiteral() (localctx IBooleanLiteralContext) {
	localctx = NewBooleanLiteralContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 78, FqlParserRULE_booleanLiteral)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(420)
		p.Match(FqlParserBooleanLiteral)
	}

	return localctx
}

// IStringLiteralContext is an interface to support dynamic dispatch.
type IStringLiteralContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsStringLiteralContext differentiates from other interfaces.
	IsStringLiteralContext()
}

type StringLiteralContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyStringLiteralContext() *StringLiteralContext {
	var p = new(StringLiteralContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = FqlParserRULE_stringLiteral
	return p
}

func (*StringLiteralContext) IsStringLiteralContext() {}

func NewStringLiteralContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *StringLiteralContext {
	var p = new(StringLiteralContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_stringLiteral

	return p
}

func (s *StringLiteralContext) GetParser() antlr.Parser { return s.parser }

func (s *StringLiteralContext) StringLiteral() antlr.TerminalNode {
	return s.GetToken(FqlParserStringLiteral, 0)
}

func (s *StringLiteralContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *StringLiteralContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *StringLiteralContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterStringLiteral(s)
	}
}

func (s *StringLiteralContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitStringLiteral(s)
	}
}

func (s *StringLiteralContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitStringLiteral(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) StringLiteral() (localctx IStringLiteralContext) {
	localctx = NewStringLiteralContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 80, FqlParserRULE_stringLiteral)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(422)
		p.Match(FqlParserStringLiteral)
	}

	return localctx
}

// IIntegerLiteralContext is an interface to support dynamic dispatch.
type IIntegerLiteralContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsIntegerLiteralContext differentiates from other interfaces.
	IsIntegerLiteralContext()
}

type IntegerLiteralContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyIntegerLiteralContext() *IntegerLiteralContext {
	var p = new(IntegerLiteralContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = FqlParserRULE_integerLiteral
	return p
}

func (*IntegerLiteralContext) IsIntegerLiteralContext() {}

func NewIntegerLiteralContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *IntegerLiteralContext {
	var p = new(IntegerLiteralContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_integerLiteral

	return p
}

func (s *IntegerLiteralContext) GetParser() antlr.Parser { return s.parser }

func (s *IntegerLiteralContext) IntegerLiteral() antlr.TerminalNode {
	return s.GetToken(FqlParserIntegerLiteral, 0)
}

func (s *IntegerLiteralContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *IntegerLiteralContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *IntegerLiteralContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterIntegerLiteral(s)
	}
}

func (s *IntegerLiteralContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitIntegerLiteral(s)
	}
}

func (s *IntegerLiteralContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitIntegerLiteral(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) IntegerLiteral() (localctx IIntegerLiteralContext) {
	localctx = NewIntegerLiteralContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 82, FqlParserRULE_integerLiteral)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(424)
		p.Match(FqlParserIntegerLiteral)
	}

	return localctx
}

// IFloatLiteralContext is an interface to support dynamic dispatch.
type IFloatLiteralContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsFloatLiteralContext differentiates from other interfaces.
	IsFloatLiteralContext()
}

type FloatLiteralContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyFloatLiteralContext() *FloatLiteralContext {
	var p = new(FloatLiteralContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = FqlParserRULE_floatLiteral
	return p
}

func (*FloatLiteralContext) IsFloatLiteralContext() {}

func NewFloatLiteralContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FloatLiteralContext {
	var p = new(FloatLiteralContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_floatLiteral

	return p
}

func (s *FloatLiteralContext) GetParser() antlr.Parser { return s.parser }

func (s *FloatLiteralContext) FloatLiteral() antlr.TerminalNode {
	return s.GetToken(FqlParserFloatLiteral, 0)
}

func (s *FloatLiteralContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FloatLiteralContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *FloatLiteralContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterFloatLiteral(s)
	}
}

func (s *FloatLiteralContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitFloatLiteral(s)
	}
}

func (s *FloatLiteralContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitFloatLiteral(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) FloatLiteral() (localctx IFloatLiteralContext) {
	localctx = NewFloatLiteralContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 84, FqlParserRULE_floatLiteral)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(426)
		p.Match(FqlParserFloatLiteral)
	}

	return localctx
}

// INoneLiteralContext is an interface to support dynamic dispatch.
type INoneLiteralContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsNoneLiteralContext differentiates from other interfaces.
	IsNoneLiteralContext()
}

type NoneLiteralContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyNoneLiteralContext() *NoneLiteralContext {
	var p = new(NoneLiteralContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = FqlParserRULE_noneLiteral
	return p
}

func (*NoneLiteralContext) IsNoneLiteralContext() {}

func NewNoneLiteralContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *NoneLiteralContext {
	var p = new(NoneLiteralContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_noneLiteral

	return p
}

func (s *NoneLiteralContext) GetParser() antlr.Parser { return s.parser }

func (s *NoneLiteralContext) Null() antlr.TerminalNode {
	return s.GetToken(FqlParserNull, 0)
}

func (s *NoneLiteralContext) None() antlr.TerminalNode {
	return s.GetToken(FqlParserNone, 0)
}

func (s *NoneLiteralContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *NoneLiteralContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *NoneLiteralContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterNoneLiteral(s)
	}
}

func (s *NoneLiteralContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitNoneLiteral(s)
	}
}

func (s *NoneLiteralContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitNoneLiteral(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) NoneLiteral() (localctx INoneLiteralContext) {
	localctx = NewNoneLiteralContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 86, FqlParserRULE_noneLiteral)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(428)
		_la = p.GetTokenStream().LA(1)

		if !(_la == FqlParserNone || _la == FqlParserNull) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}

	return localctx
}

// IArrayElementListContext is an interface to support dynamic dispatch.
type IArrayElementListContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsArrayElementListContext differentiates from other interfaces.
	IsArrayElementListContext()
}

type ArrayElementListContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyArrayElementListContext() *ArrayElementListContext {
	var p = new(ArrayElementListContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = FqlParserRULE_arrayElementList
	return p
}

func (*ArrayElementListContext) IsArrayElementListContext() {}

func NewArrayElementListContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ArrayElementListContext {
	var p = new(ArrayElementListContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_arrayElementList

	return p
}

func (s *ArrayElementListContext) GetParser() antlr.Parser { return s.parser }

func (s *ArrayElementListContext) AllExpression() []IExpressionContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IExpressionContext)(nil)).Elem())
	var tst = make([]IExpressionContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IExpressionContext)
		}
	}

	return tst
}

func (s *ArrayElementListContext) Expression(i int) IExpressionContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExpressionContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *ArrayElementListContext) AllComma() []antlr.TerminalNode {
	return s.GetTokens(FqlParserComma)
}

func (s *ArrayElementListContext) Comma(i int) antlr.TerminalNode {
	return s.GetToken(FqlParserComma, i)
}

func (s *ArrayElementListContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ArrayElementListContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ArrayElementListContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterArrayElementList(s)
	}
}

func (s *ArrayElementListContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitArrayElementList(s)
	}
}

func (s *ArrayElementListContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitArrayElementList(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) ArrayElementList() (localctx IArrayElementListContext) {
	localctx = NewArrayElementListContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 88, FqlParserRULE_arrayElementList)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(430)
		p.expression(0)
	}
	p.SetState(439)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == FqlParserComma {
		p.SetState(432)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		for ok := true; ok; ok = _la == FqlParserComma {
			{
				p.SetState(431)
				p.Match(FqlParserComma)
			}

			p.SetState(434)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)
		}
		{
			p.SetState(436)
			p.expression(0)
		}

		p.SetState(441)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}

	return localctx
}

// IPropertyAssignmentContext is an interface to support dynamic dispatch.
type IPropertyAssignmentContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsPropertyAssignmentContext differentiates from other interfaces.
	IsPropertyAssignmentContext()
}

type PropertyAssignmentContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyPropertyAssignmentContext() *PropertyAssignmentContext {
	var p = new(PropertyAssignmentContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = FqlParserRULE_propertyAssignment
	return p
}

func (*PropertyAssignmentContext) IsPropertyAssignmentContext() {}

func NewPropertyAssignmentContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *PropertyAssignmentContext {
	var p = new(PropertyAssignmentContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_propertyAssignment

	return p
}

func (s *PropertyAssignmentContext) GetParser() antlr.Parser { return s.parser }

func (s *PropertyAssignmentContext) PropertyName() IPropertyNameContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IPropertyNameContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IPropertyNameContext)
}

func (s *PropertyAssignmentContext) Colon() antlr.TerminalNode {
	return s.GetToken(FqlParserColon, 0)
}

func (s *PropertyAssignmentContext) Expression() IExpressionContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExpressionContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *PropertyAssignmentContext) ComputedPropertyName() IComputedPropertyNameContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IComputedPropertyNameContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IComputedPropertyNameContext)
}

func (s *PropertyAssignmentContext) ShorthandPropertyName() IShorthandPropertyNameContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IShorthandPropertyNameContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IShorthandPropertyNameContext)
}

func (s *PropertyAssignmentContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *PropertyAssignmentContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *PropertyAssignmentContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterPropertyAssignment(s)
	}
}

func (s *PropertyAssignmentContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitPropertyAssignment(s)
	}
}

func (s *PropertyAssignmentContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitPropertyAssignment(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) PropertyAssignment() (localctx IPropertyAssignmentContext) {
	localctx = NewPropertyAssignmentContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 90, FqlParserRULE_propertyAssignment)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(451)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 39, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(442)
			p.PropertyName()
		}
		{
			p.SetState(443)
			p.Match(FqlParserColon)
		}
		{
			p.SetState(444)
			p.expression(0)
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(446)
			p.ComputedPropertyName()
		}
		{
			p.SetState(447)
			p.Match(FqlParserColon)
		}
		{
			p.SetState(448)
			p.expression(0)
		}

	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(450)
			p.ShorthandPropertyName()
		}

	}

	return localctx
}

// IShorthandPropertyNameContext is an interface to support dynamic dispatch.
type IShorthandPropertyNameContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsShorthandPropertyNameContext differentiates from other interfaces.
	IsShorthandPropertyNameContext()
}

type ShorthandPropertyNameContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyShorthandPropertyNameContext() *ShorthandPropertyNameContext {
	var p = new(ShorthandPropertyNameContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = FqlParserRULE_shorthandPropertyName
	return p
}

func (*ShorthandPropertyNameContext) IsShorthandPropertyNameContext() {}

func NewShorthandPropertyNameContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ShorthandPropertyNameContext {
	var p = new(ShorthandPropertyNameContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_shorthandPropertyName

	return p
}

func (s *ShorthandPropertyNameContext) GetParser() antlr.Parser { return s.parser }

func (s *ShorthandPropertyNameContext) Variable() IVariableContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IVariableContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IVariableContext)
}

func (s *ShorthandPropertyNameContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ShorthandPropertyNameContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ShorthandPropertyNameContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterShorthandPropertyName(s)
	}
}

func (s *ShorthandPropertyNameContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitShorthandPropertyName(s)
	}
}

func (s *ShorthandPropertyNameContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitShorthandPropertyName(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) ShorthandPropertyName() (localctx IShorthandPropertyNameContext) {
	localctx = NewShorthandPropertyNameContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 92, FqlParserRULE_shorthandPropertyName)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(453)
		p.Variable()
	}

	return localctx
}

// IComputedPropertyNameContext is an interface to support dynamic dispatch.
type IComputedPropertyNameContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsComputedPropertyNameContext differentiates from other interfaces.
	IsComputedPropertyNameContext()
}

type ComputedPropertyNameContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyComputedPropertyNameContext() *ComputedPropertyNameContext {
	var p = new(ComputedPropertyNameContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = FqlParserRULE_computedPropertyName
	return p
}

func (*ComputedPropertyNameContext) IsComputedPropertyNameContext() {}

func NewComputedPropertyNameContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ComputedPropertyNameContext {
	var p = new(ComputedPropertyNameContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_computedPropertyName

	return p
}

func (s *ComputedPropertyNameContext) GetParser() antlr.Parser { return s.parser }

func (s *ComputedPropertyNameContext) OpenBracket() antlr.TerminalNode {
	return s.GetToken(FqlParserOpenBracket, 0)
}

func (s *ComputedPropertyNameContext) Expression() IExpressionContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExpressionContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *ComputedPropertyNameContext) CloseBracket() antlr.TerminalNode {
	return s.GetToken(FqlParserCloseBracket, 0)
}

func (s *ComputedPropertyNameContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ComputedPropertyNameContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ComputedPropertyNameContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterComputedPropertyName(s)
	}
}

func (s *ComputedPropertyNameContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitComputedPropertyName(s)
	}
}

func (s *ComputedPropertyNameContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitComputedPropertyName(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) ComputedPropertyName() (localctx IComputedPropertyNameContext) {
	localctx = NewComputedPropertyNameContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 94, FqlParserRULE_computedPropertyName)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(455)
		p.Match(FqlParserOpenBracket)
	}
	{
		p.SetState(456)
		p.expression(0)
	}
	{
		p.SetState(457)
		p.Match(FqlParserCloseBracket)
	}

	return localctx
}

// IPropertyNameContext is an interface to support dynamic dispatch.
type IPropertyNameContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsPropertyNameContext differentiates from other interfaces.
	IsPropertyNameContext()
}

type PropertyNameContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyPropertyNameContext() *PropertyNameContext {
	var p = new(PropertyNameContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = FqlParserRULE_propertyName
	return p
}

func (*PropertyNameContext) IsPropertyNameContext() {}

func NewPropertyNameContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *PropertyNameContext {
	var p = new(PropertyNameContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_propertyName

	return p
}

func (s *PropertyNameContext) GetParser() antlr.Parser { return s.parser }

func (s *PropertyNameContext) Identifier() antlr.TerminalNode {
	return s.GetToken(FqlParserIdentifier, 0)
}

func (s *PropertyNameContext) StringLiteral() IStringLiteralContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IStringLiteralContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IStringLiteralContext)
}

func (s *PropertyNameContext) Param() IParamContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IParamContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IParamContext)
}

func (s *PropertyNameContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *PropertyNameContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *PropertyNameContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterPropertyName(s)
	}
}

func (s *PropertyNameContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitPropertyName(s)
	}
}

func (s *PropertyNameContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitPropertyName(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) PropertyName() (localctx IPropertyNameContext) {
	localctx = NewPropertyNameContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 96, FqlParserRULE_propertyName)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(462)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case FqlParserIdentifier:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(459)
			p.Match(FqlParserIdentifier)
		}

	case FqlParserStringLiteral:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(460)
			p.StringLiteral()
		}

	case FqlParserParam:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(461)
			p.Param()
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}

	return localctx
}

// IExpressionGroupContext is an interface to support dynamic dispatch.
type IExpressionGroupContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsExpressionGroupContext differentiates from other interfaces.
	IsExpressionGroupContext()
}

type ExpressionGroupContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyExpressionGroupContext() *ExpressionGroupContext {
	var p = new(ExpressionGroupContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = FqlParserRULE_expressionGroup
	return p
}

func (*ExpressionGroupContext) IsExpressionGroupContext() {}

func NewExpressionGroupContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ExpressionGroupContext {
	var p = new(ExpressionGroupContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_expressionGroup

	return p
}

func (s *ExpressionGroupContext) GetParser() antlr.Parser { return s.parser }

func (s *ExpressionGroupContext) OpenParen() antlr.TerminalNode {
	return s.GetToken(FqlParserOpenParen, 0)
}

func (s *ExpressionGroupContext) Expression() IExpressionContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExpressionContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *ExpressionGroupContext) CloseParen() antlr.TerminalNode {
	return s.GetToken(FqlParserCloseParen, 0)
}

func (s *ExpressionGroupContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ExpressionGroupContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ExpressionGroupContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterExpressionGroup(s)
	}
}

func (s *ExpressionGroupContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitExpressionGroup(s)
	}
}

func (s *ExpressionGroupContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitExpressionGroup(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) ExpressionGroup() (localctx IExpressionGroupContext) {
	localctx = NewExpressionGroupContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 98, FqlParserRULE_expressionGroup)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(464)
		p.Match(FqlParserOpenParen)
	}
	{
		p.SetState(465)
		p.expression(0)
	}
	{
		p.SetState(466)
		p.Match(FqlParserCloseParen)
	}

	return localctx
}

// INamespaceIdentifierContext is an interface to support dynamic dispatch.
type INamespaceIdentifierContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsNamespaceIdentifierContext differentiates from other interfaces.
	IsNamespaceIdentifierContext()
}

type NamespaceIdentifierContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyNamespaceIdentifierContext() *NamespaceIdentifierContext {
	var p = new(NamespaceIdentifierContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = FqlParserRULE_namespaceIdentifier
	return p
}

func (*NamespaceIdentifierContext) IsNamespaceIdentifierContext() {}

func NewNamespaceIdentifierContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *NamespaceIdentifierContext {
	var p = new(NamespaceIdentifierContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_namespaceIdentifier

	return p
}

func (s *NamespaceIdentifierContext) GetParser() antlr.Parser { return s.parser }

func (s *NamespaceIdentifierContext) Namespace() INamespaceContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*INamespaceContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(INamespaceContext)
}

func (s *NamespaceIdentifierContext) Identifier() antlr.TerminalNode {
	return s.GetToken(FqlParserIdentifier, 0)
}

func (s *NamespaceIdentifierContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *NamespaceIdentifierContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *NamespaceIdentifierContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterNamespaceIdentifier(s)
	}
}

func (s *NamespaceIdentifierContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitNamespaceIdentifier(s)
	}
}

func (s *NamespaceIdentifierContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitNamespaceIdentifier(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) NamespaceIdentifier() (localctx INamespaceIdentifierContext) {
	localctx = NewNamespaceIdentifierContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 100, FqlParserRULE_namespaceIdentifier)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(468)
		p.Namespace()
	}
	{
		p.SetState(469)
		p.Match(FqlParserIdentifier)
	}

	return localctx
}

// INamespaceContext is an interface to support dynamic dispatch.
type INamespaceContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsNamespaceContext differentiates from other interfaces.
	IsNamespaceContext()
}

type NamespaceContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyNamespaceContext() *NamespaceContext {
	var p = new(NamespaceContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = FqlParserRULE_namespace
	return p
}

func (*NamespaceContext) IsNamespaceContext() {}

func NewNamespaceContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *NamespaceContext {
	var p = new(NamespaceContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_namespace

	return p
}

func (s *NamespaceContext) GetParser() antlr.Parser { return s.parser }

func (s *NamespaceContext) AllNamespaceSegment() []antlr.TerminalNode {
	return s.GetTokens(FqlParserNamespaceSegment)
}

func (s *NamespaceContext) NamespaceSegment(i int) antlr.TerminalNode {
	return s.GetToken(FqlParserNamespaceSegment, i)
}

func (s *NamespaceContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *NamespaceContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *NamespaceContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterNamespace(s)
	}
}

func (s *NamespaceContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitNamespace(s)
	}
}

func (s *NamespaceContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitNamespace(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) Namespace() (localctx INamespaceContext) {
	localctx = NewNamespaceContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 102, FqlParserRULE_namespace)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	p.SetState(474)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == FqlParserNamespaceSegment {
		{
			p.SetState(471)
			p.Match(FqlParserNamespaceSegment)
		}

		p.SetState(476)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}

	return localctx
}

// IFunctionIdentifierContext is an interface to support dynamic dispatch.
type IFunctionIdentifierContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsFunctionIdentifierContext differentiates from other interfaces.
	IsFunctionIdentifierContext()
}

type FunctionIdentifierContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyFunctionIdentifierContext() *FunctionIdentifierContext {
	var p = new(FunctionIdentifierContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = FqlParserRULE_functionIdentifier
	return p
}

func (*FunctionIdentifierContext) IsFunctionIdentifierContext() {}

func NewFunctionIdentifierContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FunctionIdentifierContext {
	var p = new(FunctionIdentifierContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_functionIdentifier

	return p
}

func (s *FunctionIdentifierContext) GetParser() antlr.Parser { return s.parser }

func (s *FunctionIdentifierContext) Identifier() antlr.TerminalNode {
	return s.GetToken(FqlParserIdentifier, 0)
}

func (s *FunctionIdentifierContext) And() antlr.TerminalNode {
	return s.GetToken(FqlParserAnd, 0)
}

func (s *FunctionIdentifierContext) Or() antlr.TerminalNode {
	return s.GetToken(FqlParserOr, 0)
}

func (s *FunctionIdentifierContext) For() antlr.TerminalNode {
	return s.GetToken(FqlParserFor, 0)
}

func (s *FunctionIdentifierContext) Return() antlr.TerminalNode {
	return s.GetToken(FqlParserReturn, 0)
}

func (s *FunctionIdentifierContext) Distinct() antlr.TerminalNode {
	return s.GetToken(FqlParserDistinct, 0)
}

func (s *FunctionIdentifierContext) Filter() antlr.TerminalNode {
	return s.GetToken(FqlParserFilter, 0)
}

func (s *FunctionIdentifierContext) Sort() antlr.TerminalNode {
	return s.GetToken(FqlParserSort, 0)
}

func (s *FunctionIdentifierContext) Limit() antlr.TerminalNode {
	return s.GetToken(FqlParserLimit, 0)
}

func (s *FunctionIdentifierContext) Let() antlr.TerminalNode {
	return s.GetToken(FqlParserLet, 0)
}

func (s *FunctionIdentifierContext) Collect() antlr.TerminalNode {
	return s.GetToken(FqlParserCollect, 0)
}

func (s *FunctionIdentifierContext) SortDirection() antlr.TerminalNode {
	return s.GetToken(FqlParserSortDirection, 0)
}

func (s *FunctionIdentifierContext) None() antlr.TerminalNode {
	return s.GetToken(FqlParserNone, 0)
}

func (s *FunctionIdentifierContext) Null() antlr.TerminalNode {
	return s.GetToken(FqlParserNull, 0)
}

func (s *FunctionIdentifierContext) BooleanLiteral() antlr.TerminalNode {
	return s.GetToken(FqlParserBooleanLiteral, 0)
}

func (s *FunctionIdentifierContext) Use() antlr.TerminalNode {
	return s.GetToken(FqlParserUse, 0)
}

func (s *FunctionIdentifierContext) Into() antlr.TerminalNode {
	return s.GetToken(FqlParserInto, 0)
}

func (s *FunctionIdentifierContext) Keep() antlr.TerminalNode {
	return s.GetToken(FqlParserKeep, 0)
}

func (s *FunctionIdentifierContext) With() antlr.TerminalNode {
	return s.GetToken(FqlParserWith, 0)
}

func (s *FunctionIdentifierContext) Count() antlr.TerminalNode {
	return s.GetToken(FqlParserCount, 0)
}

func (s *FunctionIdentifierContext) All() antlr.TerminalNode {
	return s.GetToken(FqlParserAll, 0)
}

func (s *FunctionIdentifierContext) Any() antlr.TerminalNode {
	return s.GetToken(FqlParserAny, 0)
}

func (s *FunctionIdentifierContext) Aggregate() antlr.TerminalNode {
	return s.GetToken(FqlParserAggregate, 0)
}

func (s *FunctionIdentifierContext) Like() antlr.TerminalNode {
	return s.GetToken(FqlParserLike, 0)
}

func (s *FunctionIdentifierContext) Not() antlr.TerminalNode {
	return s.GetToken(FqlParserNot, 0)
}

func (s *FunctionIdentifierContext) In() antlr.TerminalNode {
	return s.GetToken(FqlParserIn, 0)
}

func (s *FunctionIdentifierContext) Waitfor() antlr.TerminalNode {
	return s.GetToken(FqlParserWaitfor, 0)
}

func (s *FunctionIdentifierContext) Event() antlr.TerminalNode {
	return s.GetToken(FqlParserEvent, 0)
}

func (s *FunctionIdentifierContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FunctionIdentifierContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *FunctionIdentifierContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterFunctionIdentifier(s)
	}
}

func (s *FunctionIdentifierContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitFunctionIdentifier(s)
	}
}

func (s *FunctionIdentifierContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitFunctionIdentifier(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) FunctionIdentifier() (localctx IFunctionIdentifierContext) {
	localctx = NewFunctionIdentifierContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 104, FqlParserRULE_functionIdentifier)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(477)
		_la = p.GetTokenStream().LA(1)

		if !(_la == FqlParserAnd || _la == FqlParserOr || (((_la-35)&-(0x1f+1)) == 0 && ((1<<uint((_la-35)))&((1<<(FqlParserFor-35))|(1<<(FqlParserReturn-35))|(1<<(FqlParserWaitfor-35))|(1<<(FqlParserDistinct-35))|(1<<(FqlParserFilter-35))|(1<<(FqlParserSort-35))|(1<<(FqlParserLimit-35))|(1<<(FqlParserLet-35))|(1<<(FqlParserCollect-35))|(1<<(FqlParserSortDirection-35))|(1<<(FqlParserNone-35))|(1<<(FqlParserNull-35))|(1<<(FqlParserBooleanLiteral-35))|(1<<(FqlParserUse-35))|(1<<(FqlParserInto-35))|(1<<(FqlParserKeep-35))|(1<<(FqlParserWith-35))|(1<<(FqlParserCount-35))|(1<<(FqlParserAll-35))|(1<<(FqlParserAny-35))|(1<<(FqlParserAggregate-35))|(1<<(FqlParserEvent-35))|(1<<(FqlParserLike-35))|(1<<(FqlParserNot-35))|(1<<(FqlParserIn-35))|(1<<(FqlParserIdentifier-35)))) != 0)) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}

	return localctx
}

// IFunctionCallExpressionContext is an interface to support dynamic dispatch.
type IFunctionCallExpressionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsFunctionCallExpressionContext differentiates from other interfaces.
	IsFunctionCallExpressionContext()
}

type FunctionCallExpressionContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyFunctionCallExpressionContext() *FunctionCallExpressionContext {
	var p = new(FunctionCallExpressionContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = FqlParserRULE_functionCallExpression
	return p
}

func (*FunctionCallExpressionContext) IsFunctionCallExpressionContext() {}

func NewFunctionCallExpressionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FunctionCallExpressionContext {
	var p = new(FunctionCallExpressionContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_functionCallExpression

	return p
}

func (s *FunctionCallExpressionContext) GetParser() antlr.Parser { return s.parser }

func (s *FunctionCallExpressionContext) Namespace() INamespaceContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*INamespaceContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(INamespaceContext)
}

func (s *FunctionCallExpressionContext) FunctionIdentifier() IFunctionIdentifierContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IFunctionIdentifierContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IFunctionIdentifierContext)
}

func (s *FunctionCallExpressionContext) Arguments() IArgumentsContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IArgumentsContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IArgumentsContext)
}

func (s *FunctionCallExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FunctionCallExpressionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *FunctionCallExpressionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterFunctionCallExpression(s)
	}
}

func (s *FunctionCallExpressionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitFunctionCallExpression(s)
	}
}

func (s *FunctionCallExpressionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitFunctionCallExpression(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) FunctionCallExpression() (localctx IFunctionCallExpressionContext) {
	localctx = NewFunctionCallExpressionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 106, FqlParserRULE_functionCallExpression)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(479)
		p.Namespace()
	}
	{
		p.SetState(480)
		p.FunctionIdentifier()
	}
	{
		p.SetState(481)
		p.Arguments()
	}

	return localctx
}

// IMemberContext is an interface to support dynamic dispatch.
type IMemberContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsMemberContext differentiates from other interfaces.
	IsMemberContext()
}

type MemberContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyMemberContext() *MemberContext {
	var p = new(MemberContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = FqlParserRULE_member
	return p
}

func (*MemberContext) IsMemberContext() {}

func NewMemberContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *MemberContext {
	var p = new(MemberContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_member

	return p
}

func (s *MemberContext) GetParser() antlr.Parser { return s.parser }

func (s *MemberContext) Identifier() antlr.TerminalNode {
	return s.GetToken(FqlParserIdentifier, 0)
}

func (s *MemberContext) FunctionCallExpression() IFunctionCallExpressionContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IFunctionCallExpressionContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IFunctionCallExpressionContext)
}

func (s *MemberContext) Param() IParamContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IParamContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IParamContext)
}

func (s *MemberContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *MemberContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *MemberContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterMember(s)
	}
}

func (s *MemberContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitMember(s)
	}
}

func (s *MemberContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitMember(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) Member() (localctx IMemberContext) {
	localctx = NewMemberContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 108, FqlParserRULE_member)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(486)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 42, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(483)
			p.Match(FqlParserIdentifier)
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(484)
			p.FunctionCallExpression()
		}

	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(485)
			p.Param()
		}

	}

	return localctx
}

// IMemberPathContext is an interface to support dynamic dispatch.
type IMemberPathContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsMemberPathContext differentiates from other interfaces.
	IsMemberPathContext()
}

type MemberPathContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyMemberPathContext() *MemberPathContext {
	var p = new(MemberPathContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = FqlParserRULE_memberPath
	return p
}

func (*MemberPathContext) IsMemberPathContext() {}

func NewMemberPathContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *MemberPathContext {
	var p = new(MemberPathContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_memberPath

	return p
}

func (s *MemberPathContext) GetParser() antlr.Parser { return s.parser }

func (s *MemberPathContext) AllDot() []antlr.TerminalNode {
	return s.GetTokens(FqlParserDot)
}

func (s *MemberPathContext) Dot(i int) antlr.TerminalNode {
	return s.GetToken(FqlParserDot, i)
}

func (s *MemberPathContext) AllPropertyName() []IPropertyNameContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IPropertyNameContext)(nil)).Elem())
	var tst = make([]IPropertyNameContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IPropertyNameContext)
		}
	}

	return tst
}

func (s *MemberPathContext) PropertyName(i int) IPropertyNameContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IPropertyNameContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IPropertyNameContext)
}

func (s *MemberPathContext) AllComputedPropertyName() []IComputedPropertyNameContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IComputedPropertyNameContext)(nil)).Elem())
	var tst = make([]IComputedPropertyNameContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IComputedPropertyNameContext)
		}
	}

	return tst
}

func (s *MemberPathContext) ComputedPropertyName(i int) IComputedPropertyNameContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IComputedPropertyNameContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IComputedPropertyNameContext)
}

func (s *MemberPathContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *MemberPathContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *MemberPathContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterMemberPath(s)
	}
}

func (s *MemberPathContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitMemberPath(s)
	}
}

func (s *MemberPathContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitMemberPath(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) MemberPath() (localctx IMemberPathContext) {
	localctx = NewMemberPathContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 110, FqlParserRULE_memberPath)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	var _alt int

	p.SetState(527)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case FqlParserDot:
		p.EnterOuterAlt(localctx, 1)
		p.SetState(496)
		p.GetErrorHandler().Sync(p)
		_alt = 1
		for ok := true; ok; ok = _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
			switch _alt {
			case 1:
				{
					p.SetState(488)
					p.Match(FqlParserDot)
				}
				{
					p.SetState(489)
					p.PropertyName()
				}
				p.SetState(493)
				p.GetErrorHandler().Sync(p)
				_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 43, p.GetParserRuleContext())

				for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
					if _alt == 1 {
						{
							p.SetState(490)
							p.ComputedPropertyName()
						}

					}
					p.SetState(495)
					p.GetErrorHandler().Sync(p)
					_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 43, p.GetParserRuleContext())
				}

			default:
				panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
			}

			p.SetState(498)
			p.GetErrorHandler().Sync(p)
			_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 44, p.GetParserRuleContext())
		}

	case FqlParserOpenBracket:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(500)
			p.ComputedPropertyName()
		}
		p.SetState(511)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 46, p.GetParserRuleContext())

		for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
			if _alt == 1 {
				{
					p.SetState(501)
					p.Match(FqlParserDot)
				}
				{
					p.SetState(502)
					p.PropertyName()
				}
				p.SetState(506)
				p.GetErrorHandler().Sync(p)
				_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 45, p.GetParserRuleContext())

				for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
					if _alt == 1 {
						{
							p.SetState(503)
							p.ComputedPropertyName()
						}

					}
					p.SetState(508)
					p.GetErrorHandler().Sync(p)
					_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 45, p.GetParserRuleContext())
				}

			}
			p.SetState(513)
			p.GetErrorHandler().Sync(p)
			_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 46, p.GetParserRuleContext())
		}
		p.SetState(524)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 48, p.GetParserRuleContext())

		for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
			if _alt == 1 {
				{
					p.SetState(514)
					p.ComputedPropertyName()
				}
				p.SetState(519)
				p.GetErrorHandler().Sync(p)
				_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 47, p.GetParserRuleContext())

				for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
					if _alt == 1 {
						{
							p.SetState(515)
							p.Match(FqlParserDot)
						}
						{
							p.SetState(516)
							p.PropertyName()
						}

					}
					p.SetState(521)
					p.GetErrorHandler().Sync(p)
					_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 47, p.GetParserRuleContext())
				}

			}
			p.SetState(526)
			p.GetErrorHandler().Sync(p)
			_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 48, p.GetParserRuleContext())
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}

	return localctx
}

// IMemberExpressionContext is an interface to support dynamic dispatch.
type IMemberExpressionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsMemberExpressionContext differentiates from other interfaces.
	IsMemberExpressionContext()
}

type MemberExpressionContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyMemberExpressionContext() *MemberExpressionContext {
	var p = new(MemberExpressionContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = FqlParserRULE_memberExpression
	return p
}

func (*MemberExpressionContext) IsMemberExpressionContext() {}

func NewMemberExpressionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *MemberExpressionContext {
	var p = new(MemberExpressionContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_memberExpression

	return p
}

func (s *MemberExpressionContext) GetParser() antlr.Parser { return s.parser }

func (s *MemberExpressionContext) Member() IMemberContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IMemberContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IMemberContext)
}

func (s *MemberExpressionContext) MemberPath() IMemberPathContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IMemberPathContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IMemberPathContext)
}

func (s *MemberExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *MemberExpressionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *MemberExpressionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterMemberExpression(s)
	}
}

func (s *MemberExpressionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitMemberExpression(s)
	}
}

func (s *MemberExpressionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitMemberExpression(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) MemberExpression() (localctx IMemberExpressionContext) {
	localctx = NewMemberExpressionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 112, FqlParserRULE_memberExpression)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(529)
		p.Member()
	}
	{
		p.SetState(530)
		p.MemberPath()
	}

	return localctx
}

// IArgumentsContext is an interface to support dynamic dispatch.
type IArgumentsContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsArgumentsContext differentiates from other interfaces.
	IsArgumentsContext()
}

type ArgumentsContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyArgumentsContext() *ArgumentsContext {
	var p = new(ArgumentsContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = FqlParserRULE_arguments
	return p
}

func (*ArgumentsContext) IsArgumentsContext() {}

func NewArgumentsContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ArgumentsContext {
	var p = new(ArgumentsContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_arguments

	return p
}

func (s *ArgumentsContext) GetParser() antlr.Parser { return s.parser }

func (s *ArgumentsContext) OpenParen() antlr.TerminalNode {
	return s.GetToken(FqlParserOpenParen, 0)
}

func (s *ArgumentsContext) CloseParen() antlr.TerminalNode {
	return s.GetToken(FqlParserCloseParen, 0)
}

func (s *ArgumentsContext) AllExpression() []IExpressionContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IExpressionContext)(nil)).Elem())
	var tst = make([]IExpressionContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IExpressionContext)
		}
	}

	return tst
}

func (s *ArgumentsContext) Expression(i int) IExpressionContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExpressionContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *ArgumentsContext) AllComma() []antlr.TerminalNode {
	return s.GetTokens(FqlParserComma)
}

func (s *ArgumentsContext) Comma(i int) antlr.TerminalNode {
	return s.GetToken(FqlParserComma, i)
}

func (s *ArgumentsContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ArgumentsContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ArgumentsContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterArguments(s)
	}
}

func (s *ArgumentsContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitArguments(s)
	}
}

func (s *ArgumentsContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitArguments(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) Arguments() (localctx IArgumentsContext) {
	localctx = NewArgumentsContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 114, FqlParserRULE_arguments)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(532)
		p.Match(FqlParserOpenParen)
	}
	p.SetState(541)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if (((_la-9)&-(0x1f+1)) == 0 && ((1<<uint((_la-9)))&((1<<(FqlParserOpenBracket-9))|(1<<(FqlParserOpenParen-9))|(1<<(FqlParserOpenBrace-9))|(1<<(FqlParserPlus-9))|(1<<(FqlParserMinus-9))|(1<<(FqlParserAnd-9))|(1<<(FqlParserOr-9))|(1<<(FqlParserFor-9))|(1<<(FqlParserReturn-9))|(1<<(FqlParserWaitfor-9))|(1<<(FqlParserDistinct-9))|(1<<(FqlParserFilter-9))|(1<<(FqlParserSort-9)))) != 0) || (((_la-41)&-(0x1f+1)) == 0 && ((1<<uint((_la-41)))&((1<<(FqlParserLimit-41))|(1<<(FqlParserLet-41))|(1<<(FqlParserCollect-41))|(1<<(FqlParserSortDirection-41))|(1<<(FqlParserNone-41))|(1<<(FqlParserNull-41))|(1<<(FqlParserBooleanLiteral-41))|(1<<(FqlParserUse-41))|(1<<(FqlParserInto-41))|(1<<(FqlParserKeep-41))|(1<<(FqlParserWith-41))|(1<<(FqlParserCount-41))|(1<<(FqlParserAll-41))|(1<<(FqlParserAny-41))|(1<<(FqlParserAggregate-41))|(1<<(FqlParserEvent-41))|(1<<(FqlParserLike-41))|(1<<(FqlParserNot-41))|(1<<(FqlParserIn-41))|(1<<(FqlParserParam-41))|(1<<(FqlParserIdentifier-41))|(1<<(FqlParserStringLiteral-41))|(1<<(FqlParserIntegerLiteral-41))|(1<<(FqlParserFloatLiteral-41))|(1<<(FqlParserNamespaceSegment-41)))) != 0) {
		{
			p.SetState(533)
			p.expression(0)
		}
		p.SetState(538)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		for _la == FqlParserComma {
			{
				p.SetState(534)
				p.Match(FqlParserComma)
			}
			{
				p.SetState(535)
				p.expression(0)
			}

			p.SetState(540)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)
		}

	}
	{
		p.SetState(543)
		p.Match(FqlParserCloseParen)
	}

	return localctx
}

// IExpressionContext is an interface to support dynamic dispatch.
type IExpressionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsExpressionContext differentiates from other interfaces.
	IsExpressionContext()
}

type ExpressionContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyExpressionContext() *ExpressionContext {
	var p = new(ExpressionContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = FqlParserRULE_expression
	return p
}

func (*ExpressionContext) IsExpressionContext() {}

func NewExpressionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ExpressionContext {
	var p = new(ExpressionContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_expression

	return p
}

func (s *ExpressionContext) GetParser() antlr.Parser { return s.parser }

func (s *ExpressionContext) UnaryOperator() IUnaryOperatorContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IUnaryOperatorContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IUnaryOperatorContext)
}

func (s *ExpressionContext) AllExpression() []IExpressionContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IExpressionContext)(nil)).Elem())
	var tst = make([]IExpressionContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IExpressionContext)
		}
	}

	return tst
}

func (s *ExpressionContext) Expression(i int) IExpressionContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExpressionContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *ExpressionContext) FunctionCallExpression() IFunctionCallExpressionContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IFunctionCallExpressionContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IFunctionCallExpressionContext)
}

func (s *ExpressionContext) ExpressionGroup() IExpressionGroupContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExpressionGroupContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IExpressionGroupContext)
}

func (s *ExpressionContext) RangeOperator() IRangeOperatorContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IRangeOperatorContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IRangeOperatorContext)
}

func (s *ExpressionContext) StringLiteral() IStringLiteralContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IStringLiteralContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IStringLiteralContext)
}

func (s *ExpressionContext) IntegerLiteral() IIntegerLiteralContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IIntegerLiteralContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IIntegerLiteralContext)
}

func (s *ExpressionContext) FloatLiteral() IFloatLiteralContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IFloatLiteralContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IFloatLiteralContext)
}

func (s *ExpressionContext) BooleanLiteral() IBooleanLiteralContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IBooleanLiteralContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IBooleanLiteralContext)
}

func (s *ExpressionContext) ArrayLiteral() IArrayLiteralContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IArrayLiteralContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IArrayLiteralContext)
}

func (s *ExpressionContext) ObjectLiteral() IObjectLiteralContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IObjectLiteralContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IObjectLiteralContext)
}

func (s *ExpressionContext) Variable() IVariableContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IVariableContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IVariableContext)
}

func (s *ExpressionContext) MemberExpression() IMemberExpressionContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IMemberExpressionContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IMemberExpressionContext)
}

func (s *ExpressionContext) NoneLiteral() INoneLiteralContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*INoneLiteralContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(INoneLiteralContext)
}

func (s *ExpressionContext) Param() IParamContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IParamContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IParamContext)
}

func (s *ExpressionContext) MultiplicativeOperator() IMultiplicativeOperatorContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IMultiplicativeOperatorContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IMultiplicativeOperatorContext)
}

func (s *ExpressionContext) AdditiveOperator() IAdditiveOperatorContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IAdditiveOperatorContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IAdditiveOperatorContext)
}

func (s *ExpressionContext) ArrayOperator() IArrayOperatorContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IArrayOperatorContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IArrayOperatorContext)
}

func (s *ExpressionContext) InOperator() IInOperatorContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IInOperatorContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IInOperatorContext)
}

func (s *ExpressionContext) EqualityOperator() IEqualityOperatorContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IEqualityOperatorContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IEqualityOperatorContext)
}

func (s *ExpressionContext) RegexpOperator() IRegexpOperatorContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IRegexpOperatorContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IRegexpOperatorContext)
}

func (s *ExpressionContext) LogicalAndOperator() ILogicalAndOperatorContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ILogicalAndOperatorContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ILogicalAndOperatorContext)
}

func (s *ExpressionContext) LogicalOrOperator() ILogicalOrOperatorContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ILogicalOrOperatorContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ILogicalOrOperatorContext)
}

func (s *ExpressionContext) QuestionMark() antlr.TerminalNode {
	return s.GetToken(FqlParserQuestionMark, 0)
}

func (s *ExpressionContext) Colon() antlr.TerminalNode {
	return s.GetToken(FqlParserColon, 0)
}

func (s *ExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ExpressionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ExpressionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterExpression(s)
	}
}

func (s *ExpressionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitExpression(s)
	}
}

func (s *ExpressionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitExpression(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) Expression() (localctx IExpressionContext) {
	return p.expression(0)
}

func (p *FqlParser) expression(_p int) (localctx IExpressionContext) {
	var _parentctx antlr.ParserRuleContext = p.GetParserRuleContext()
	_parentState := p.GetState()
	localctx = NewExpressionContext(p, p.GetParserRuleContext(), _parentState)
	var _prevctx IExpressionContext = localctx
	var _ antlr.ParserRuleContext = _prevctx // TODO: To prevent unused variable warning.
	_startState := 116
	p.EnterRecursionRule(localctx, 116, FqlParserRULE_expression, _p)
	var _la int

	defer func() {
		p.UnrollRecursionContexts(_parentctx)
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	var _alt int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(562)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 52, p.GetParserRuleContext()) {
	case 1:
		{
			p.SetState(546)
			p.UnaryOperator()
		}
		{
			p.SetState(547)
			p.expression(23)
		}

	case 2:
		{
			p.SetState(549)
			p.FunctionCallExpression()
		}

	case 3:
		{
			p.SetState(550)
			p.ExpressionGroup()
		}

	case 4:
		{
			p.SetState(551)
			p.RangeOperator()
		}

	case 5:
		{
			p.SetState(552)
			p.StringLiteral()
		}

	case 6:
		{
			p.SetState(553)
			p.IntegerLiteral()
		}

	case 7:
		{
			p.SetState(554)
			p.FloatLiteral()
		}

	case 8:
		{
			p.SetState(555)
			p.BooleanLiteral()
		}

	case 9:
		{
			p.SetState(556)
			p.ArrayLiteral()
		}

	case 10:
		{
			p.SetState(557)
			p.ObjectLiteral()
		}

	case 11:
		{
			p.SetState(558)
			p.Variable()
		}

	case 12:
		{
			p.SetState(559)
			p.MemberExpression()
		}

	case 13:
		{
			p.SetState(560)
			p.NoneLiteral()
		}

	case 14:
		{
			p.SetState(561)
			p.Param()
		}

	}
	p.GetParserRuleContext().SetStop(p.GetTokenStream().LT(-1))
	p.SetState(609)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 56, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			if p.GetParseListeners() != nil {
				p.TriggerExitRuleEvent()
			}
			_prevctx = localctx
			p.SetState(607)
			p.GetErrorHandler().Sync(p)
			switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 55, p.GetParserRuleContext()) {
			case 1:
				localctx = NewExpressionContext(p, _parentctx, _parentState)
				p.PushNewRecursionContext(localctx, _startState, FqlParserRULE_expression)
				p.SetState(564)

				if !(p.Precpred(p.GetParserRuleContext(), 22)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 22)", ""))
				}
				{
					p.SetState(565)
					p.MultiplicativeOperator()
				}
				{
					p.SetState(566)
					p.expression(23)
				}

			case 2:
				localctx = NewExpressionContext(p, _parentctx, _parentState)
				p.PushNewRecursionContext(localctx, _startState, FqlParserRULE_expression)
				p.SetState(568)

				if !(p.Precpred(p.GetParserRuleContext(), 21)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 21)", ""))
				}
				{
					p.SetState(569)
					p.AdditiveOperator()
				}
				{
					p.SetState(570)
					p.expression(22)
				}

			case 3:
				localctx = NewExpressionContext(p, _parentctx, _parentState)
				p.PushNewRecursionContext(localctx, _startState, FqlParserRULE_expression)
				p.SetState(572)

				if !(p.Precpred(p.GetParserRuleContext(), 18)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 18)", ""))
				}
				{
					p.SetState(573)
					p.ArrayOperator()
				}
				p.SetState(576)
				p.GetErrorHandler().Sync(p)

				switch p.GetTokenStream().LA(1) {
				case FqlParserNot, FqlParserIn:
					{
						p.SetState(574)
						p.InOperator()
					}

				case FqlParserGt, FqlParserLt, FqlParserEq, FqlParserGte, FqlParserLte, FqlParserNeq:
					{
						p.SetState(575)
						p.EqualityOperator()
					}

				default:
					panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
				}
				{
					p.SetState(578)
					p.expression(19)
				}

			case 4:
				localctx = NewExpressionContext(p, _parentctx, _parentState)
				p.PushNewRecursionContext(localctx, _startState, FqlParserRULE_expression)
				p.SetState(580)

				if !(p.Precpred(p.GetParserRuleContext(), 17)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 17)", ""))
				}
				{
					p.SetState(581)
					p.InOperator()
				}
				{
					p.SetState(582)
					p.expression(18)
				}

			case 5:
				localctx = NewExpressionContext(p, _parentctx, _parentState)
				p.PushNewRecursionContext(localctx, _startState, FqlParserRULE_expression)
				p.SetState(584)

				if !(p.Precpred(p.GetParserRuleContext(), 16)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 16)", ""))
				}
				{
					p.SetState(585)
					p.EqualityOperator()
				}
				{
					p.SetState(586)
					p.expression(17)
				}

			case 6:
				localctx = NewExpressionContext(p, _parentctx, _parentState)
				p.PushNewRecursionContext(localctx, _startState, FqlParserRULE_expression)
				p.SetState(588)

				if !(p.Precpred(p.GetParserRuleContext(), 15)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 15)", ""))
				}
				{
					p.SetState(589)
					p.RegexpOperator()
				}
				{
					p.SetState(590)
					p.expression(16)
				}

			case 7:
				localctx = NewExpressionContext(p, _parentctx, _parentState)
				p.PushNewRecursionContext(localctx, _startState, FqlParserRULE_expression)
				p.SetState(592)

				if !(p.Precpred(p.GetParserRuleContext(), 14)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 14)", ""))
				}
				{
					p.SetState(593)
					p.LogicalAndOperator()
				}
				{
					p.SetState(594)
					p.expression(15)
				}

			case 8:
				localctx = NewExpressionContext(p, _parentctx, _parentState)
				p.PushNewRecursionContext(localctx, _startState, FqlParserRULE_expression)
				p.SetState(596)

				if !(p.Precpred(p.GetParserRuleContext(), 13)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 13)", ""))
				}
				{
					p.SetState(597)
					p.LogicalOrOperator()
				}
				{
					p.SetState(598)
					p.expression(14)
				}

			case 9:
				localctx = NewExpressionContext(p, _parentctx, _parentState)
				p.PushNewRecursionContext(localctx, _startState, FqlParserRULE_expression)
				p.SetState(600)

				if !(p.Precpred(p.GetParserRuleContext(), 12)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 12)", ""))
				}
				{
					p.SetState(601)
					p.Match(FqlParserQuestionMark)
				}
				p.SetState(603)
				p.GetErrorHandler().Sync(p)
				_la = p.GetTokenStream().LA(1)

				if (((_la-9)&-(0x1f+1)) == 0 && ((1<<uint((_la-9)))&((1<<(FqlParserOpenBracket-9))|(1<<(FqlParserOpenParen-9))|(1<<(FqlParserOpenBrace-9))|(1<<(FqlParserPlus-9))|(1<<(FqlParserMinus-9))|(1<<(FqlParserAnd-9))|(1<<(FqlParserOr-9))|(1<<(FqlParserFor-9))|(1<<(FqlParserReturn-9))|(1<<(FqlParserWaitfor-9))|(1<<(FqlParserDistinct-9))|(1<<(FqlParserFilter-9))|(1<<(FqlParserSort-9)))) != 0) || (((_la-41)&-(0x1f+1)) == 0 && ((1<<uint((_la-41)))&((1<<(FqlParserLimit-41))|(1<<(FqlParserLet-41))|(1<<(FqlParserCollect-41))|(1<<(FqlParserSortDirection-41))|(1<<(FqlParserNone-41))|(1<<(FqlParserNull-41))|(1<<(FqlParserBooleanLiteral-41))|(1<<(FqlParserUse-41))|(1<<(FqlParserInto-41))|(1<<(FqlParserKeep-41))|(1<<(FqlParserWith-41))|(1<<(FqlParserCount-41))|(1<<(FqlParserAll-41))|(1<<(FqlParserAny-41))|(1<<(FqlParserAggregate-41))|(1<<(FqlParserEvent-41))|(1<<(FqlParserLike-41))|(1<<(FqlParserNot-41))|(1<<(FqlParserIn-41))|(1<<(FqlParserParam-41))|(1<<(FqlParserIdentifier-41))|(1<<(FqlParserStringLiteral-41))|(1<<(FqlParserIntegerLiteral-41))|(1<<(FqlParserFloatLiteral-41))|(1<<(FqlParserNamespaceSegment-41)))) != 0) {
					{
						p.SetState(602)
						p.expression(0)
					}

				}
				{
					p.SetState(605)
					p.Match(FqlParserColon)
				}
				{
					p.SetState(606)
					p.expression(13)
				}

			}

		}
		p.SetState(611)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 56, p.GetParserRuleContext())
	}

	return localctx
}

// IForTernaryExpressionContext is an interface to support dynamic dispatch.
type IForTernaryExpressionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsForTernaryExpressionContext differentiates from other interfaces.
	IsForTernaryExpressionContext()
}

type ForTernaryExpressionContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyForTernaryExpressionContext() *ForTernaryExpressionContext {
	var p = new(ForTernaryExpressionContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = FqlParserRULE_forTernaryExpression
	return p
}

func (*ForTernaryExpressionContext) IsForTernaryExpressionContext() {}

func NewForTernaryExpressionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ForTernaryExpressionContext {
	var p = new(ForTernaryExpressionContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_forTernaryExpression

	return p
}

func (s *ForTernaryExpressionContext) GetParser() antlr.Parser { return s.parser }

func (s *ForTernaryExpressionContext) AllExpression() []IExpressionContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IExpressionContext)(nil)).Elem())
	var tst = make([]IExpressionContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IExpressionContext)
		}
	}

	return tst
}

func (s *ForTernaryExpressionContext) Expression(i int) IExpressionContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExpressionContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *ForTernaryExpressionContext) QuestionMark() antlr.TerminalNode {
	return s.GetToken(FqlParserQuestionMark, 0)
}

func (s *ForTernaryExpressionContext) Colon() antlr.TerminalNode {
	return s.GetToken(FqlParserColon, 0)
}

func (s *ForTernaryExpressionContext) AllOpenParen() []antlr.TerminalNode {
	return s.GetTokens(FqlParserOpenParen)
}

func (s *ForTernaryExpressionContext) OpenParen(i int) antlr.TerminalNode {
	return s.GetToken(FqlParserOpenParen, i)
}

func (s *ForTernaryExpressionContext) AllForExpression() []IForExpressionContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IForExpressionContext)(nil)).Elem())
	var tst = make([]IForExpressionContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IForExpressionContext)
		}
	}

	return tst
}

func (s *ForTernaryExpressionContext) ForExpression(i int) IForExpressionContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IForExpressionContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IForExpressionContext)
}

func (s *ForTernaryExpressionContext) AllCloseParen() []antlr.TerminalNode {
	return s.GetTokens(FqlParserCloseParen)
}

func (s *ForTernaryExpressionContext) CloseParen(i int) antlr.TerminalNode {
	return s.GetToken(FqlParserCloseParen, i)
}

func (s *ForTernaryExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ForTernaryExpressionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ForTernaryExpressionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterForTernaryExpression(s)
	}
}

func (s *ForTernaryExpressionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitForTernaryExpression(s)
	}
}

func (s *ForTernaryExpressionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitForTernaryExpression(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) ForTernaryExpression() (localctx IForTernaryExpressionContext) {
	localctx = NewForTernaryExpressionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 118, FqlParserRULE_forTernaryExpression)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(640)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 58, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(612)
			p.expression(0)
		}
		{
			p.SetState(613)
			p.Match(FqlParserQuestionMark)
		}
		p.SetState(615)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		if (((_la-9)&-(0x1f+1)) == 0 && ((1<<uint((_la-9)))&((1<<(FqlParserOpenBracket-9))|(1<<(FqlParserOpenParen-9))|(1<<(FqlParserOpenBrace-9))|(1<<(FqlParserPlus-9))|(1<<(FqlParserMinus-9))|(1<<(FqlParserAnd-9))|(1<<(FqlParserOr-9))|(1<<(FqlParserFor-9))|(1<<(FqlParserReturn-9))|(1<<(FqlParserWaitfor-9))|(1<<(FqlParserDistinct-9))|(1<<(FqlParserFilter-9))|(1<<(FqlParserSort-9)))) != 0) || (((_la-41)&-(0x1f+1)) == 0 && ((1<<uint((_la-41)))&((1<<(FqlParserLimit-41))|(1<<(FqlParserLet-41))|(1<<(FqlParserCollect-41))|(1<<(FqlParserSortDirection-41))|(1<<(FqlParserNone-41))|(1<<(FqlParserNull-41))|(1<<(FqlParserBooleanLiteral-41))|(1<<(FqlParserUse-41))|(1<<(FqlParserInto-41))|(1<<(FqlParserKeep-41))|(1<<(FqlParserWith-41))|(1<<(FqlParserCount-41))|(1<<(FqlParserAll-41))|(1<<(FqlParserAny-41))|(1<<(FqlParserAggregate-41))|(1<<(FqlParserEvent-41))|(1<<(FqlParserLike-41))|(1<<(FqlParserNot-41))|(1<<(FqlParserIn-41))|(1<<(FqlParserParam-41))|(1<<(FqlParserIdentifier-41))|(1<<(FqlParserStringLiteral-41))|(1<<(FqlParserIntegerLiteral-41))|(1<<(FqlParserFloatLiteral-41))|(1<<(FqlParserNamespaceSegment-41)))) != 0) {
			{
				p.SetState(614)
				p.expression(0)
			}

		}
		{
			p.SetState(617)
			p.Match(FqlParserColon)
		}
		{
			p.SetState(618)
			p.Match(FqlParserOpenParen)
		}
		{
			p.SetState(619)
			p.ForExpression()
		}
		{
			p.SetState(620)
			p.Match(FqlParserCloseParen)
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(622)
			p.expression(0)
		}
		{
			p.SetState(623)
			p.Match(FqlParserQuestionMark)
		}
		{
			p.SetState(624)
			p.Match(FqlParserOpenParen)
		}
		{
			p.SetState(625)
			p.ForExpression()
		}
		{
			p.SetState(626)
			p.Match(FqlParserCloseParen)
		}
		{
			p.SetState(627)
			p.Match(FqlParserColon)
		}
		{
			p.SetState(628)
			p.expression(0)
		}

	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(630)
			p.expression(0)
		}
		{
			p.SetState(631)
			p.Match(FqlParserQuestionMark)
		}
		{
			p.SetState(632)
			p.Match(FqlParserOpenParen)
		}
		{
			p.SetState(633)
			p.ForExpression()
		}
		{
			p.SetState(634)
			p.Match(FqlParserCloseParen)
		}
		{
			p.SetState(635)
			p.Match(FqlParserColon)
		}
		{
			p.SetState(636)
			p.Match(FqlParserOpenParen)
		}
		{
			p.SetState(637)
			p.ForExpression()
		}
		{
			p.SetState(638)
			p.Match(FqlParserCloseParen)
		}

	}

	return localctx
}

// IArrayOperatorContext is an interface to support dynamic dispatch.
type IArrayOperatorContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsArrayOperatorContext differentiates from other interfaces.
	IsArrayOperatorContext()
}

type ArrayOperatorContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyArrayOperatorContext() *ArrayOperatorContext {
	var p = new(ArrayOperatorContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = FqlParserRULE_arrayOperator
	return p
}

func (*ArrayOperatorContext) IsArrayOperatorContext() {}

func NewArrayOperatorContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ArrayOperatorContext {
	var p = new(ArrayOperatorContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_arrayOperator

	return p
}

func (s *ArrayOperatorContext) GetParser() antlr.Parser { return s.parser }

func (s *ArrayOperatorContext) All() antlr.TerminalNode {
	return s.GetToken(FqlParserAll, 0)
}

func (s *ArrayOperatorContext) Any() antlr.TerminalNode {
	return s.GetToken(FqlParserAny, 0)
}

func (s *ArrayOperatorContext) None() antlr.TerminalNode {
	return s.GetToken(FqlParserNone, 0)
}

func (s *ArrayOperatorContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ArrayOperatorContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ArrayOperatorContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterArrayOperator(s)
	}
}

func (s *ArrayOperatorContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitArrayOperator(s)
	}
}

func (s *ArrayOperatorContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitArrayOperator(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) ArrayOperator() (localctx IArrayOperatorContext) {
	localctx = NewArrayOperatorContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 120, FqlParserRULE_arrayOperator)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(642)
		_la = p.GetTokenStream().LA(1)

		if !(((_la-45)&-(0x1f+1)) == 0 && ((1<<uint((_la-45)))&((1<<(FqlParserNone-45))|(1<<(FqlParserAll-45))|(1<<(FqlParserAny-45)))) != 0) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}

	return localctx
}

// IInOperatorContext is an interface to support dynamic dispatch.
type IInOperatorContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsInOperatorContext differentiates from other interfaces.
	IsInOperatorContext()
}

type InOperatorContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyInOperatorContext() *InOperatorContext {
	var p = new(InOperatorContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = FqlParserRULE_inOperator
	return p
}

func (*InOperatorContext) IsInOperatorContext() {}

func NewInOperatorContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *InOperatorContext {
	var p = new(InOperatorContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_inOperator

	return p
}

func (s *InOperatorContext) GetParser() antlr.Parser { return s.parser }

func (s *InOperatorContext) In() antlr.TerminalNode {
	return s.GetToken(FqlParserIn, 0)
}

func (s *InOperatorContext) Not() antlr.TerminalNode {
	return s.GetToken(FqlParserNot, 0)
}

func (s *InOperatorContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *InOperatorContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *InOperatorContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterInOperator(s)
	}
}

func (s *InOperatorContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitInOperator(s)
	}
}

func (s *InOperatorContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitInOperator(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) InOperator() (localctx IInOperatorContext) {
	localctx = NewInOperatorContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 122, FqlParserRULE_inOperator)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(647)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case FqlParserIn:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(644)
			p.Match(FqlParserIn)
		}

	case FqlParserNot:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(645)
			p.Match(FqlParserNot)
		}
		{
			p.SetState(646)
			p.Match(FqlParserIn)
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}

	return localctx
}

// IEqualityOperatorContext is an interface to support dynamic dispatch.
type IEqualityOperatorContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsEqualityOperatorContext differentiates from other interfaces.
	IsEqualityOperatorContext()
}

type EqualityOperatorContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyEqualityOperatorContext() *EqualityOperatorContext {
	var p = new(EqualityOperatorContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = FqlParserRULE_equalityOperator
	return p
}

func (*EqualityOperatorContext) IsEqualityOperatorContext() {}

func NewEqualityOperatorContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *EqualityOperatorContext {
	var p = new(EqualityOperatorContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_equalityOperator

	return p
}

func (s *EqualityOperatorContext) GetParser() antlr.Parser { return s.parser }

func (s *EqualityOperatorContext) Gt() antlr.TerminalNode {
	return s.GetToken(FqlParserGt, 0)
}

func (s *EqualityOperatorContext) Lt() antlr.TerminalNode {
	return s.GetToken(FqlParserLt, 0)
}

func (s *EqualityOperatorContext) Eq() antlr.TerminalNode {
	return s.GetToken(FqlParserEq, 0)
}

func (s *EqualityOperatorContext) Gte() antlr.TerminalNode {
	return s.GetToken(FqlParserGte, 0)
}

func (s *EqualityOperatorContext) Lte() antlr.TerminalNode {
	return s.GetToken(FqlParserLte, 0)
}

func (s *EqualityOperatorContext) Neq() antlr.TerminalNode {
	return s.GetToken(FqlParserNeq, 0)
}

func (s *EqualityOperatorContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *EqualityOperatorContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *EqualityOperatorContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterEqualityOperator(s)
	}
}

func (s *EqualityOperatorContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitEqualityOperator(s)
	}
}

func (s *EqualityOperatorContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitEqualityOperator(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) EqualityOperator() (localctx IEqualityOperatorContext) {
	localctx = NewEqualityOperatorContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 124, FqlParserRULE_equalityOperator)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(649)
		_la = p.GetTokenStream().LA(1)

		if !(((_la)&-(0x1f+1)) == 0 && ((1<<uint(_la))&((1<<FqlParserGt)|(1<<FqlParserLt)|(1<<FqlParserEq)|(1<<FqlParserGte)|(1<<FqlParserLte)|(1<<FqlParserNeq))) != 0) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}

	return localctx
}

// IRegexpOperatorContext is an interface to support dynamic dispatch.
type IRegexpOperatorContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsRegexpOperatorContext differentiates from other interfaces.
	IsRegexpOperatorContext()
}

type RegexpOperatorContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyRegexpOperatorContext() *RegexpOperatorContext {
	var p = new(RegexpOperatorContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = FqlParserRULE_regexpOperator
	return p
}

func (*RegexpOperatorContext) IsRegexpOperatorContext() {}

func NewRegexpOperatorContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *RegexpOperatorContext {
	var p = new(RegexpOperatorContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_regexpOperator

	return p
}

func (s *RegexpOperatorContext) GetParser() antlr.Parser { return s.parser }

func (s *RegexpOperatorContext) RegexMatch() antlr.TerminalNode {
	return s.GetToken(FqlParserRegexMatch, 0)
}

func (s *RegexpOperatorContext) RegexNotMatch() antlr.TerminalNode {
	return s.GetToken(FqlParserRegexNotMatch, 0)
}

func (s *RegexpOperatorContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *RegexpOperatorContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *RegexpOperatorContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterRegexpOperator(s)
	}
}

func (s *RegexpOperatorContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitRegexpOperator(s)
	}
}

func (s *RegexpOperatorContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitRegexpOperator(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) RegexpOperator() (localctx IRegexpOperatorContext) {
	localctx = NewRegexpOperatorContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 126, FqlParserRULE_regexpOperator)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(651)
		_la = p.GetTokenStream().LA(1)

		if !(_la == FqlParserRegexNotMatch || _la == FqlParserRegexMatch) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}

	return localctx
}

// ILogicalAndOperatorContext is an interface to support dynamic dispatch.
type ILogicalAndOperatorContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsLogicalAndOperatorContext differentiates from other interfaces.
	IsLogicalAndOperatorContext()
}

type LogicalAndOperatorContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyLogicalAndOperatorContext() *LogicalAndOperatorContext {
	var p = new(LogicalAndOperatorContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = FqlParserRULE_logicalAndOperator
	return p
}

func (*LogicalAndOperatorContext) IsLogicalAndOperatorContext() {}

func NewLogicalAndOperatorContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *LogicalAndOperatorContext {
	var p = new(LogicalAndOperatorContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_logicalAndOperator

	return p
}

func (s *LogicalAndOperatorContext) GetParser() antlr.Parser { return s.parser }

func (s *LogicalAndOperatorContext) And() antlr.TerminalNode {
	return s.GetToken(FqlParserAnd, 0)
}

func (s *LogicalAndOperatorContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *LogicalAndOperatorContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *LogicalAndOperatorContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterLogicalAndOperator(s)
	}
}

func (s *LogicalAndOperatorContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitLogicalAndOperator(s)
	}
}

func (s *LogicalAndOperatorContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitLogicalAndOperator(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) LogicalAndOperator() (localctx ILogicalAndOperatorContext) {
	localctx = NewLogicalAndOperatorContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 128, FqlParserRULE_logicalAndOperator)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(653)
		p.Match(FqlParserAnd)
	}

	return localctx
}

// ILogicalOrOperatorContext is an interface to support dynamic dispatch.
type ILogicalOrOperatorContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsLogicalOrOperatorContext differentiates from other interfaces.
	IsLogicalOrOperatorContext()
}

type LogicalOrOperatorContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyLogicalOrOperatorContext() *LogicalOrOperatorContext {
	var p = new(LogicalOrOperatorContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = FqlParserRULE_logicalOrOperator
	return p
}

func (*LogicalOrOperatorContext) IsLogicalOrOperatorContext() {}

func NewLogicalOrOperatorContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *LogicalOrOperatorContext {
	var p = new(LogicalOrOperatorContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_logicalOrOperator

	return p
}

func (s *LogicalOrOperatorContext) GetParser() antlr.Parser { return s.parser }

func (s *LogicalOrOperatorContext) Or() antlr.TerminalNode {
	return s.GetToken(FqlParserOr, 0)
}

func (s *LogicalOrOperatorContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *LogicalOrOperatorContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *LogicalOrOperatorContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterLogicalOrOperator(s)
	}
}

func (s *LogicalOrOperatorContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitLogicalOrOperator(s)
	}
}

func (s *LogicalOrOperatorContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitLogicalOrOperator(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) LogicalOrOperator() (localctx ILogicalOrOperatorContext) {
	localctx = NewLogicalOrOperatorContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 130, FqlParserRULE_logicalOrOperator)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(655)
		p.Match(FqlParserOr)
	}

	return localctx
}

// IMultiplicativeOperatorContext is an interface to support dynamic dispatch.
type IMultiplicativeOperatorContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsMultiplicativeOperatorContext differentiates from other interfaces.
	IsMultiplicativeOperatorContext()
}

type MultiplicativeOperatorContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyMultiplicativeOperatorContext() *MultiplicativeOperatorContext {
	var p = new(MultiplicativeOperatorContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = FqlParserRULE_multiplicativeOperator
	return p
}

func (*MultiplicativeOperatorContext) IsMultiplicativeOperatorContext() {}

func NewMultiplicativeOperatorContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *MultiplicativeOperatorContext {
	var p = new(MultiplicativeOperatorContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_multiplicativeOperator

	return p
}

func (s *MultiplicativeOperatorContext) GetParser() antlr.Parser { return s.parser }

func (s *MultiplicativeOperatorContext) Multi() antlr.TerminalNode {
	return s.GetToken(FqlParserMulti, 0)
}

func (s *MultiplicativeOperatorContext) Div() antlr.TerminalNode {
	return s.GetToken(FqlParserDiv, 0)
}

func (s *MultiplicativeOperatorContext) Mod() antlr.TerminalNode {
	return s.GetToken(FqlParserMod, 0)
}

func (s *MultiplicativeOperatorContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *MultiplicativeOperatorContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *MultiplicativeOperatorContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterMultiplicativeOperator(s)
	}
}

func (s *MultiplicativeOperatorContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitMultiplicativeOperator(s)
	}
}

func (s *MultiplicativeOperatorContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitMultiplicativeOperator(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) MultiplicativeOperator() (localctx IMultiplicativeOperatorContext) {
	localctx = NewMultiplicativeOperatorContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 132, FqlParserRULE_multiplicativeOperator)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(657)
		_la = p.GetTokenStream().LA(1)

		if !(((_la)&-(0x1f+1)) == 0 && ((1<<uint(_la))&((1<<FqlParserMulti)|(1<<FqlParserDiv)|(1<<FqlParserMod))) != 0) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}

	return localctx
}

// IAdditiveOperatorContext is an interface to support dynamic dispatch.
type IAdditiveOperatorContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsAdditiveOperatorContext differentiates from other interfaces.
	IsAdditiveOperatorContext()
}

type AdditiveOperatorContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyAdditiveOperatorContext() *AdditiveOperatorContext {
	var p = new(AdditiveOperatorContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = FqlParserRULE_additiveOperator
	return p
}

func (*AdditiveOperatorContext) IsAdditiveOperatorContext() {}

func NewAdditiveOperatorContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *AdditiveOperatorContext {
	var p = new(AdditiveOperatorContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_additiveOperator

	return p
}

func (s *AdditiveOperatorContext) GetParser() antlr.Parser { return s.parser }

func (s *AdditiveOperatorContext) Plus() antlr.TerminalNode {
	return s.GetToken(FqlParserPlus, 0)
}

func (s *AdditiveOperatorContext) Minus() antlr.TerminalNode {
	return s.GetToken(FqlParserMinus, 0)
}

func (s *AdditiveOperatorContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *AdditiveOperatorContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *AdditiveOperatorContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterAdditiveOperator(s)
	}
}

func (s *AdditiveOperatorContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitAdditiveOperator(s)
	}
}

func (s *AdditiveOperatorContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitAdditiveOperator(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) AdditiveOperator() (localctx IAdditiveOperatorContext) {
	localctx = NewAdditiveOperatorContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 134, FqlParserRULE_additiveOperator)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(659)
		_la = p.GetTokenStream().LA(1)

		if !(_la == FqlParserPlus || _la == FqlParserMinus) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}

	return localctx
}

// IUnaryOperatorContext is an interface to support dynamic dispatch.
type IUnaryOperatorContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsUnaryOperatorContext differentiates from other interfaces.
	IsUnaryOperatorContext()
}

type UnaryOperatorContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyUnaryOperatorContext() *UnaryOperatorContext {
	var p = new(UnaryOperatorContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = FqlParserRULE_unaryOperator
	return p
}

func (*UnaryOperatorContext) IsUnaryOperatorContext() {}

func NewUnaryOperatorContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *UnaryOperatorContext {
	var p = new(UnaryOperatorContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_unaryOperator

	return p
}

func (s *UnaryOperatorContext) GetParser() antlr.Parser { return s.parser }

func (s *UnaryOperatorContext) Not() antlr.TerminalNode {
	return s.GetToken(FqlParserNot, 0)
}

func (s *UnaryOperatorContext) Plus() antlr.TerminalNode {
	return s.GetToken(FqlParserPlus, 0)
}

func (s *UnaryOperatorContext) Minus() antlr.TerminalNode {
	return s.GetToken(FqlParserMinus, 0)
}

func (s *UnaryOperatorContext) Like() antlr.TerminalNode {
	return s.GetToken(FqlParserLike, 0)
}

func (s *UnaryOperatorContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *UnaryOperatorContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *UnaryOperatorContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterUnaryOperator(s)
	}
}

func (s *UnaryOperatorContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitUnaryOperator(s)
	}
}

func (s *UnaryOperatorContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitUnaryOperator(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) UnaryOperator() (localctx IUnaryOperatorContext) {
	localctx = NewUnaryOperatorContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 136, FqlParserRULE_unaryOperator)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(661)
		_la = p.GetTokenStream().LA(1)

		if !(_la == FqlParserPlus || _la == FqlParserMinus || _la == FqlParserLike || _la == FqlParserNot) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}

	return localctx
}

func (p *FqlParser) Sempred(localctx antlr.RuleContext, ruleIndex, predIndex int) bool {
	switch ruleIndex {
	case 58:
		var t *ExpressionContext = nil
		if localctx != nil {
			t = localctx.(*ExpressionContext)
		}
		return p.Expression_Sempred(t, predIndex)

	default:
		panic("No predicate with index: " + fmt.Sprint(ruleIndex))
	}
}

func (p *FqlParser) Expression_Sempred(localctx antlr.RuleContext, predIndex int) bool {
	switch predIndex {
	case 0:
		return p.Precpred(p.GetParserRuleContext(), 22)

	case 1:
		return p.Precpred(p.GetParserRuleContext(), 21)

	case 2:
		return p.Precpred(p.GetParserRuleContext(), 18)

	case 3:
		return p.Precpred(p.GetParserRuleContext(), 17)

	case 4:
		return p.Precpred(p.GetParserRuleContext(), 16)

	case 5:
		return p.Precpred(p.GetParserRuleContext(), 15)

	case 6:
		return p.Precpred(p.GetParserRuleContext(), 14)

	case 7:
		return p.Precpred(p.GetParserRuleContext(), 13)

	case 8:
		return p.Precpred(p.GetParserRuleContext(), 12)

	default:
		panic("No predicate with index: " + fmt.Sprint(predIndex))
	}
}
