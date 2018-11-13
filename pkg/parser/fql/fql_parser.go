// Code generated from antlr/FqlParser.g4 by ANTLR 4.7.1. DO NOT EDIT.

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
	3, 24715, 42794, 33075, 47597, 16764, 15335, 30598, 22884, 3, 64, 545,
	4, 2, 9, 2, 4, 3, 9, 3, 4, 4, 9, 4, 4, 5, 9, 5, 4, 6, 9, 6, 4, 7, 9, 7,
	4, 8, 9, 8, 4, 9, 9, 9, 4, 10, 9, 10, 4, 11, 9, 11, 4, 12, 9, 12, 4, 13,
	9, 13, 4, 14, 9, 14, 4, 15, 9, 15, 4, 16, 9, 16, 4, 17, 9, 17, 4, 18, 9,
	18, 4, 19, 9, 19, 4, 20, 9, 20, 4, 21, 9, 21, 4, 22, 9, 22, 4, 23, 9, 23,
	4, 24, 9, 24, 4, 25, 9, 25, 4, 26, 9, 26, 4, 27, 9, 27, 4, 28, 9, 28, 4,
	29, 9, 29, 4, 30, 9, 30, 4, 31, 9, 31, 4, 32, 9, 32, 4, 33, 9, 33, 4, 34,
	9, 34, 4, 35, 9, 35, 4, 36, 9, 36, 4, 37, 9, 37, 4, 38, 9, 38, 4, 39, 9,
	39, 4, 40, 9, 40, 4, 41, 9, 41, 4, 42, 9, 42, 4, 43, 9, 43, 4, 44, 9, 44,
	4, 45, 9, 45, 4, 46, 9, 46, 4, 47, 9, 47, 4, 48, 9, 48, 4, 49, 9, 49, 4,
	50, 9, 50, 4, 51, 9, 51, 4, 52, 9, 52, 4, 53, 9, 53, 4, 54, 9, 54, 3, 2,
	3, 2, 3, 3, 7, 3, 112, 10, 3, 12, 3, 14, 3, 115, 11, 3, 3, 3, 3, 3, 3,
	4, 3, 4, 5, 4, 121, 10, 4, 3, 5, 3, 5, 5, 5, 125, 10, 5, 3, 6, 3, 6, 5,
	6, 129, 10, 6, 3, 6, 3, 6, 3, 6, 5, 6, 134, 10, 6, 3, 6, 3, 6, 3, 6, 3,
	6, 3, 6, 3, 6, 5, 6, 142, 10, 6, 3, 7, 3, 7, 3, 7, 3, 7, 5, 7, 148, 10,
	7, 3, 7, 3, 7, 3, 7, 7, 7, 153, 10, 7, 12, 7, 14, 7, 156, 11, 7, 3, 7,
	3, 7, 3, 8, 3, 8, 3, 9, 3, 9, 3, 10, 3, 10, 3, 10, 3, 10, 3, 10, 3, 10,
	3, 10, 5, 10, 171, 10, 10, 3, 11, 3, 11, 3, 11, 3, 11, 5, 11, 177, 10,
	11, 3, 12, 3, 12, 5, 12, 181, 10, 12, 3, 13, 3, 13, 5, 13, 185, 10, 13,
	3, 14, 3, 14, 5, 14, 189, 10, 14, 3, 15, 3, 15, 3, 15, 3, 16, 3, 16, 3,
	16, 3, 16, 5, 16, 198, 10, 16, 3, 17, 3, 17, 5, 17, 202, 10, 17, 3, 18,
	3, 18, 3, 18, 3, 18, 7, 18, 208, 10, 18, 12, 18, 14, 18, 211, 11, 18, 3,
	19, 3, 19, 5, 19, 215, 10, 19, 3, 20, 3, 20, 3, 20, 3, 20, 3, 20, 3, 20,
	3, 20, 3, 20, 3, 20, 3, 20, 3, 20, 3, 20, 3, 20, 3, 20, 3, 20, 3, 20, 3,
	20, 3, 20, 5, 20, 235, 10, 20, 3, 21, 3, 21, 3, 21, 3, 21, 3, 22, 3, 22,
	3, 22, 7, 22, 244, 10, 22, 12, 22, 14, 22, 247, 11, 22, 3, 23, 3, 23, 3,
	23, 3, 23, 7, 23, 253, 10, 23, 12, 23, 14, 23, 256, 11, 23, 3, 24, 3, 24,
	3, 24, 3, 24, 3, 25, 3, 25, 3, 25, 3, 25, 3, 25, 3, 25, 5, 25, 268, 10,
	25, 5, 25, 270, 10, 25, 3, 26, 3, 26, 3, 26, 3, 26, 3, 26, 3, 27, 3, 27,
	3, 27, 3, 27, 3, 27, 3, 27, 3, 27, 3, 27, 3, 27, 3, 27, 3, 27, 3, 27, 3,
	27, 3, 27, 3, 27, 5, 27, 292, 10, 27, 3, 28, 3, 28, 3, 28, 3, 29, 3, 29,
	3, 30, 3, 30, 3, 30, 5, 30, 302, 10, 30, 3, 30, 3, 30, 3, 30, 3, 30, 5,
	30, 308, 10, 30, 3, 31, 3, 31, 5, 31, 312, 10, 31, 3, 31, 3, 31, 3, 32,
	3, 32, 3, 32, 3, 32, 7, 32, 320, 10, 32, 12, 32, 14, 32, 323, 11, 32, 5,
	32, 325, 10, 32, 3, 32, 5, 32, 328, 10, 32, 3, 32, 3, 32, 3, 33, 3, 33,
	3, 34, 3, 34, 3, 35, 3, 35, 3, 36, 3, 36, 3, 37, 3, 37, 3, 38, 3, 38, 6,
	38, 344, 10, 38, 13, 38, 14, 38, 345, 3, 38, 7, 38, 349, 10, 38, 12, 38,
	14, 38, 352, 11, 38, 3, 39, 3, 39, 3, 39, 3, 39, 3, 39, 3, 39, 3, 39, 3,
	39, 3, 39, 5, 39, 363, 10, 39, 3, 40, 3, 40, 3, 40, 3, 40, 7, 40, 369,
	10, 40, 12, 40, 14, 40, 372, 11, 40, 6, 40, 374, 10, 40, 13, 40, 14, 40,
	375, 3, 40, 3, 40, 3, 40, 3, 40, 3, 40, 7, 40, 383, 10, 40, 12, 40, 14,
	40, 386, 11, 40, 7, 40, 388, 10, 40, 12, 40, 14, 40, 391, 11, 40, 3, 40,
	3, 40, 3, 40, 7, 40, 396, 10, 40, 12, 40, 14, 40, 399, 11, 40, 7, 40, 401,
	10, 40, 12, 40, 14, 40, 404, 11, 40, 5, 40, 406, 10, 40, 3, 41, 3, 41,
	3, 42, 3, 42, 3, 42, 3, 42, 3, 43, 3, 43, 5, 43, 416, 10, 43, 3, 44, 3,
	44, 3, 44, 7, 44, 421, 10, 44, 12, 44, 14, 44, 424, 11, 44, 3, 45, 3, 45,
	3, 45, 3, 46, 3, 46, 3, 46, 3, 46, 7, 46, 433, 10, 46, 12, 46, 14, 46,
	436, 11, 46, 5, 46, 438, 10, 46, 3, 46, 3, 46, 3, 47, 3, 47, 3, 47, 3,
	47, 3, 47, 3, 47, 3, 47, 3, 47, 3, 47, 3, 47, 3, 47, 3, 47, 3, 47, 3, 47,
	3, 47, 3, 47, 3, 47, 3, 47, 3, 47, 3, 47, 5, 47, 462, 10, 47, 3, 47, 3,
	47, 3, 47, 3, 47, 3, 47, 3, 47, 3, 47, 3, 47, 3, 47, 3, 47, 3, 47, 3, 47,
	3, 47, 3, 47, 3, 47, 3, 47, 5, 47, 480, 10, 47, 3, 47, 3, 47, 3, 47, 3,
	47, 3, 47, 3, 47, 3, 47, 3, 47, 3, 47, 5, 47, 491, 10, 47, 3, 47, 3, 47,
	7, 47, 495, 10, 47, 12, 47, 14, 47, 498, 11, 47, 3, 48, 3, 48, 3, 48, 5,
	48, 503, 10, 48, 3, 48, 3, 48, 3, 48, 3, 48, 3, 48, 3, 48, 3, 48, 3, 48,
	3, 48, 3, 48, 3, 48, 3, 48, 3, 48, 3, 48, 3, 48, 3, 48, 3, 48, 3, 48, 3,
	48, 3, 48, 3, 48, 3, 48, 3, 48, 5, 48, 528, 10, 48, 3, 49, 3, 49, 3, 50,
	3, 50, 3, 50, 5, 50, 535, 10, 50, 3, 51, 3, 51, 3, 52, 3, 52, 3, 53, 3,
	53, 3, 54, 3, 54, 3, 54, 2, 3, 92, 55, 2, 4, 6, 8, 10, 12, 14, 16, 18,
	20, 22, 24, 26, 28, 30, 32, 34, 36, 38, 40, 42, 44, 46, 48, 50, 52, 54,
	56, 58, 60, 62, 64, 66, 68, 70, 72, 74, 76, 78, 80, 82, 84, 86, 88, 90,
	92, 94, 96, 98, 100, 102, 104, 106, 2, 9, 3, 2, 61, 62, 3, 2, 46, 47, 4,
	2, 46, 46, 53, 54, 3, 2, 17, 22, 3, 2, 30, 31, 4, 2, 23, 24, 27, 29, 4,
	2, 23, 24, 56, 57, 2, 575, 2, 108, 3, 2, 2, 2, 4, 113, 3, 2, 2, 2, 6, 120,
	3, 2, 2, 2, 8, 124, 3, 2, 2, 2, 10, 141, 3, 2, 2, 2, 12, 143, 3, 2, 2,
	2, 14, 159, 3, 2, 2, 2, 16, 161, 3, 2, 2, 2, 18, 170, 3, 2, 2, 2, 20, 176,
	3, 2, 2, 2, 22, 180, 3, 2, 2, 2, 24, 184, 3, 2, 2, 2, 26, 188, 3, 2, 2,
	2, 28, 190, 3, 2, 2, 2, 30, 193, 3, 2, 2, 2, 32, 201, 3, 2, 2, 2, 34, 203,
	3, 2, 2, 2, 36, 212, 3, 2, 2, 2, 38, 234, 3, 2, 2, 2, 40, 236, 3, 2, 2,
	2, 42, 240, 3, 2, 2, 2, 44, 248, 3, 2, 2, 2, 46, 257, 3, 2, 2, 2, 48, 269,
	3, 2, 2, 2, 50, 271, 3, 2, 2, 2, 52, 291, 3, 2, 2, 2, 54, 293, 3, 2, 2,
	2, 56, 296, 3, 2, 2, 2, 58, 301, 3, 2, 2, 2, 60, 309, 3, 2, 2, 2, 62, 315,
	3, 2, 2, 2, 64, 331, 3, 2, 2, 2, 66, 333, 3, 2, 2, 2, 68, 335, 3, 2, 2,
	2, 70, 337, 3, 2, 2, 2, 72, 339, 3, 2, 2, 2, 74, 341, 3, 2, 2, 2, 76, 362,
	3, 2, 2, 2, 78, 405, 3, 2, 2, 2, 80, 407, 3, 2, 2, 2, 82, 409, 3, 2, 2,
	2, 84, 415, 3, 2, 2, 2, 86, 417, 3, 2, 2, 2, 88, 425, 3, 2, 2, 2, 90, 428,
	3, 2, 2, 2, 92, 461, 3, 2, 2, 2, 94, 527, 3, 2, 2, 2, 96, 529, 3, 2, 2,
	2, 98, 534, 3, 2, 2, 2, 100, 536, 3, 2, 2, 2, 102, 538, 3, 2, 2, 2, 104,
	540, 3, 2, 2, 2, 106, 542, 3, 2, 2, 2, 108, 109, 5, 4, 3, 2, 109, 3, 3,
	2, 2, 2, 110, 112, 5, 6, 4, 2, 111, 110, 3, 2, 2, 2, 112, 115, 3, 2, 2,
	2, 113, 111, 3, 2, 2, 2, 113, 114, 3, 2, 2, 2, 114, 116, 3, 2, 2, 2, 115,
	113, 3, 2, 2, 2, 116, 117, 5, 8, 5, 2, 117, 5, 3, 2, 2, 2, 118, 121, 5,
	88, 45, 2, 119, 121, 5, 52, 27, 2, 120, 118, 3, 2, 2, 2, 120, 119, 3, 2,
	2, 2, 121, 7, 3, 2, 2, 2, 122, 125, 5, 10, 6, 2, 123, 125, 5, 12, 7, 2,
	124, 122, 3, 2, 2, 2, 124, 123, 3, 2, 2, 2, 125, 9, 3, 2, 2, 2, 126, 128,
	7, 38, 2, 2, 127, 129, 7, 39, 2, 2, 128, 127, 3, 2, 2, 2, 128, 129, 3,
	2, 2, 2, 129, 130, 3, 2, 2, 2, 130, 142, 5, 92, 47, 2, 131, 133, 7, 38,
	2, 2, 132, 134, 7, 39, 2, 2, 133, 132, 3, 2, 2, 2, 133, 134, 3, 2, 2, 2,
	134, 135, 3, 2, 2, 2, 135, 136, 7, 13, 2, 2, 136, 137, 5, 12, 7, 2, 137,
	138, 7, 14, 2, 2, 138, 142, 3, 2, 2, 2, 139, 140, 7, 38, 2, 2, 140, 142,
	5, 94, 48, 2, 141, 126, 3, 2, 2, 2, 141, 131, 3, 2, 2, 2, 141, 139, 3,
	2, 2, 2, 142, 11, 3, 2, 2, 2, 143, 144, 7, 37, 2, 2, 144, 147, 5, 14, 8,
	2, 145, 146, 7, 10, 2, 2, 146, 148, 5, 16, 9, 2, 147, 145, 3, 2, 2, 2,
	147, 148, 3, 2, 2, 2, 148, 149, 3, 2, 2, 2, 149, 150, 7, 58, 2, 2, 150,
	154, 5, 18, 10, 2, 151, 153, 5, 24, 13, 2, 152, 151, 3, 2, 2, 2, 153, 156,
	3, 2, 2, 2, 154, 152, 3, 2, 2, 2, 154, 155, 3, 2, 2, 2, 155, 157, 3, 2,
	2, 2, 156, 154, 3, 2, 2, 2, 157, 158, 5, 26, 14, 2, 158, 13, 3, 2, 2, 2,
	159, 160, 7, 60, 2, 2, 160, 15, 3, 2, 2, 2, 161, 162, 7, 60, 2, 2, 162,
	17, 3, 2, 2, 2, 163, 171, 5, 88, 45, 2, 164, 171, 5, 60, 31, 2, 165, 171,
	5, 62, 32, 2, 166, 171, 5, 56, 29, 2, 167, 171, 5, 78, 40, 2, 168, 171,
	5, 58, 30, 2, 169, 171, 5, 54, 28, 2, 170, 163, 3, 2, 2, 2, 170, 164, 3,
	2, 2, 2, 170, 165, 3, 2, 2, 2, 170, 166, 3, 2, 2, 2, 170, 167, 3, 2, 2,
	2, 170, 168, 3, 2, 2, 2, 170, 169, 3, 2, 2, 2, 171, 19, 3, 2, 2, 2, 172,
	177, 5, 30, 16, 2, 173, 177, 5, 34, 18, 2, 174, 177, 5, 28, 15, 2, 175,
	177, 5, 38, 20, 2, 176, 172, 3, 2, 2, 2, 176, 173, 3, 2, 2, 2, 176, 174,
	3, 2, 2, 2, 176, 175, 3, 2, 2, 2, 177, 21, 3, 2, 2, 2, 178, 181, 5, 52,
	27, 2, 179, 181, 5, 88, 45, 2, 180, 178, 3, 2, 2, 2, 180, 179, 3, 2, 2,
	2, 181, 23, 3, 2, 2, 2, 182, 185, 5, 22, 12, 2, 183, 185, 5, 20, 11, 2,
	184, 182, 3, 2, 2, 2, 184, 183, 3, 2, 2, 2, 185, 25, 3, 2, 2, 2, 186, 189,
	5, 10, 6, 2, 187, 189, 5, 12, 7, 2, 188, 186, 3, 2, 2, 2, 188, 187, 3,
	2, 2, 2, 189, 27, 3, 2, 2, 2, 190, 191, 7, 40, 2, 2, 191, 192, 5, 92, 47,
	2, 192, 29, 3, 2, 2, 2, 193, 194, 7, 42, 2, 2, 194, 197, 5, 32, 17, 2,
	195, 196, 7, 10, 2, 2, 196, 198, 5, 32, 17, 2, 197, 195, 3, 2, 2, 2, 197,
	198, 3, 2, 2, 2, 198, 31, 3, 2, 2, 2, 199, 202, 7, 63, 2, 2, 200, 202,
	5, 54, 28, 2, 201, 199, 3, 2, 2, 2, 201, 200, 3, 2, 2, 2, 202, 33, 3, 2,
	2, 2, 203, 204, 7, 41, 2, 2, 204, 209, 5, 36, 19, 2, 205, 206, 7, 10, 2,
	2, 206, 208, 5, 36, 19, 2, 207, 205, 3, 2, 2, 2, 208, 211, 3, 2, 2, 2,
	209, 207, 3, 2, 2, 2, 209, 210, 3, 2, 2, 2, 210, 35, 3, 2, 2, 2, 211, 209,
	3, 2, 2, 2, 212, 214, 5, 92, 47, 2, 213, 215, 7, 45, 2, 2, 214, 213, 3,
	2, 2, 2, 214, 215, 3, 2, 2, 2, 215, 37, 3, 2, 2, 2, 216, 217, 7, 44, 2,
	2, 217, 235, 5, 50, 26, 2, 218, 219, 7, 44, 2, 2, 219, 235, 5, 44, 23,
	2, 220, 221, 7, 44, 2, 2, 221, 222, 5, 42, 22, 2, 222, 223, 5, 44, 23,
	2, 223, 235, 3, 2, 2, 2, 224, 225, 7, 44, 2, 2, 225, 226, 5, 42, 22, 2,
	226, 227, 5, 48, 25, 2, 227, 235, 3, 2, 2, 2, 228, 229, 7, 44, 2, 2, 229,
	230, 5, 42, 22, 2, 230, 231, 5, 50, 26, 2, 231, 235, 3, 2, 2, 2, 232, 233,
	7, 44, 2, 2, 233, 235, 5, 42, 22, 2, 234, 216, 3, 2, 2, 2, 234, 218, 3,
	2, 2, 2, 234, 220, 3, 2, 2, 2, 234, 224, 3, 2, 2, 2, 234, 228, 3, 2, 2,
	2, 234, 232, 3, 2, 2, 2, 235, 39, 3, 2, 2, 2, 236, 237, 7, 60, 2, 2, 237,
	238, 7, 33, 2, 2, 238, 239, 5, 92, 47, 2, 239, 41, 3, 2, 2, 2, 240, 245,
	5, 40, 21, 2, 241, 242, 7, 10, 2, 2, 242, 244, 5, 40, 21, 2, 243, 241,
	3, 2, 2, 2, 244, 247, 3, 2, 2, 2, 245, 243, 3, 2, 2, 2, 245, 246, 3, 2,
	2, 2, 246, 43, 3, 2, 2, 2, 247, 245, 3, 2, 2, 2, 248, 249, 7, 55, 2, 2,
	249, 254, 5, 46, 24, 2, 250, 251, 7, 10, 2, 2, 251, 253, 5, 46, 24, 2,
	252, 250, 3, 2, 2, 2, 253, 256, 3, 2, 2, 2, 254, 252, 3, 2, 2, 2, 254,
	255, 3, 2, 2, 2, 255, 45, 3, 2, 2, 2, 256, 254, 3, 2, 2, 2, 257, 258, 7,
	60, 2, 2, 258, 259, 7, 33, 2, 2, 259, 260, 5, 88, 45, 2, 260, 47, 3, 2,
	2, 2, 261, 262, 7, 49, 2, 2, 262, 270, 5, 40, 21, 2, 263, 264, 7, 49, 2,
	2, 264, 267, 7, 60, 2, 2, 265, 266, 7, 50, 2, 2, 266, 268, 7, 60, 2, 2,
	267, 265, 3, 2, 2, 2, 267, 268, 3, 2, 2, 2, 268, 270, 3, 2, 2, 2, 269,
	261, 3, 2, 2, 2, 269, 263, 3, 2, 2, 2, 270, 49, 3, 2, 2, 2, 271, 272, 7,
	51, 2, 2, 272, 273, 7, 52, 2, 2, 273, 274, 7, 49, 2, 2, 274, 275, 7, 60,
	2, 2, 275, 51, 3, 2, 2, 2, 276, 277, 7, 43, 2, 2, 277, 278, 7, 60, 2, 2,
	278, 279, 7, 33, 2, 2, 279, 292, 5, 92, 47, 2, 280, 281, 7, 43, 2, 2, 281,
	282, 7, 60, 2, 2, 282, 283, 7, 33, 2, 2, 283, 284, 7, 13, 2, 2, 284, 285,
	5, 12, 7, 2, 285, 286, 7, 14, 2, 2, 286, 292, 3, 2, 2, 2, 287, 288, 7,
	43, 2, 2, 288, 289, 7, 60, 2, 2, 289, 290, 7, 33, 2, 2, 290, 292, 5, 94,
	48, 2, 291, 276, 3, 2, 2, 2, 291, 280, 3, 2, 2, 2, 291, 287, 3, 2, 2, 2,
	292, 53, 3, 2, 2, 2, 293, 294, 7, 59, 2, 2, 294, 295, 7, 60, 2, 2, 295,
	55, 3, 2, 2, 2, 296, 297, 7, 60, 2, 2, 297, 57, 3, 2, 2, 2, 298, 302, 5,
	68, 35, 2, 299, 302, 5, 56, 29, 2, 300, 302, 5, 54, 28, 2, 301, 298, 3,
	2, 2, 2, 301, 299, 3, 2, 2, 2, 301, 300, 3, 2, 2, 2, 302, 303, 3, 2, 2,
	2, 303, 307, 7, 32, 2, 2, 304, 308, 5, 68, 35, 2, 305, 308, 5, 56, 29,
	2, 306, 308, 5, 54, 28, 2, 307, 304, 3, 2, 2, 2, 307, 305, 3, 2, 2, 2,
	307, 306, 3, 2, 2, 2, 308, 59, 3, 2, 2, 2, 309, 311, 7, 11, 2, 2, 310,
	312, 5, 74, 38, 2, 311, 310, 3, 2, 2, 2, 311, 312, 3, 2, 2, 2, 312, 313,
	3, 2, 2, 2, 313, 314, 7, 12, 2, 2, 314, 61, 3, 2, 2, 2, 315, 324, 7, 15,
	2, 2, 316, 321, 5, 76, 39, 2, 317, 318, 7, 10, 2, 2, 318, 320, 5, 76, 39,
	2, 319, 317, 3, 2, 2, 2, 320, 323, 3, 2, 2, 2, 321, 319, 3, 2, 2, 2, 321,
	322, 3, 2, 2, 2, 322, 325, 3, 2, 2, 2, 323, 321, 3, 2, 2, 2, 324, 316,
	3, 2, 2, 2, 324, 325, 3, 2, 2, 2, 325, 327, 3, 2, 2, 2, 326, 328, 7, 10,
	2, 2, 327, 326, 3, 2, 2, 2, 327, 328, 3, 2, 2, 2, 328, 329, 3, 2, 2, 2,
	329, 330, 7, 16, 2, 2, 330, 63, 3, 2, 2, 2, 331, 332, 7, 48, 2, 2, 332,
	65, 3, 2, 2, 2, 333, 334, 9, 2, 2, 2, 334, 67, 3, 2, 2, 2, 335, 336, 7,
	63, 2, 2, 336, 69, 3, 2, 2, 2, 337, 338, 7, 64, 2, 2, 338, 71, 3, 2, 2,
	2, 339, 340, 9, 3, 2, 2, 340, 73, 3, 2, 2, 2, 341, 350, 5, 92, 47, 2, 342,
	344, 7, 10, 2, 2, 343, 342, 3, 2, 2, 2, 344, 345, 3, 2, 2, 2, 345, 343,
	3, 2, 2, 2, 345, 346, 3, 2, 2, 2, 346, 347, 3, 2, 2, 2, 347, 349, 5, 92,
	47, 2, 348, 343, 3, 2, 2, 2, 349, 352, 3, 2, 2, 2, 350, 348, 3, 2, 2, 2,
	350, 351, 3, 2, 2, 2, 351, 75, 3, 2, 2, 2, 352, 350, 3, 2, 2, 2, 353, 354,
	5, 84, 43, 2, 354, 355, 7, 7, 2, 2, 355, 356, 5, 92, 47, 2, 356, 363, 3,
	2, 2, 2, 357, 358, 5, 82, 42, 2, 358, 359, 7, 7, 2, 2, 359, 360, 5, 92,
	47, 2, 360, 363, 3, 2, 2, 2, 361, 363, 5, 80, 41, 2, 362, 353, 3, 2, 2,
	2, 362, 357, 3, 2, 2, 2, 362, 361, 3, 2, 2, 2, 363, 77, 3, 2, 2, 2, 364,
	373, 7, 60, 2, 2, 365, 366, 7, 9, 2, 2, 366, 370, 5, 84, 43, 2, 367, 369,
	5, 82, 42, 2, 368, 367, 3, 2, 2, 2, 369, 372, 3, 2, 2, 2, 370, 368, 3,
	2, 2, 2, 370, 371, 3, 2, 2, 2, 371, 374, 3, 2, 2, 2, 372, 370, 3, 2, 2,
	2, 373, 365, 3, 2, 2, 2, 374, 375, 3, 2, 2, 2, 375, 373, 3, 2, 2, 2, 375,
	376, 3, 2, 2, 2, 376, 406, 3, 2, 2, 2, 377, 378, 7, 60, 2, 2, 378, 389,
	5, 82, 42, 2, 379, 380, 7, 9, 2, 2, 380, 384, 5, 84, 43, 2, 381, 383, 5,
	82, 42, 2, 382, 381, 3, 2, 2, 2, 383, 386, 3, 2, 2, 2, 384, 382, 3, 2,
	2, 2, 384, 385, 3, 2, 2, 2, 385, 388, 3, 2, 2, 2, 386, 384, 3, 2, 2, 2,
	387, 379, 3, 2, 2, 2, 388, 391, 3, 2, 2, 2, 389, 387, 3, 2, 2, 2, 389,
	390, 3, 2, 2, 2, 390, 402, 3, 2, 2, 2, 391, 389, 3, 2, 2, 2, 392, 397,
	5, 82, 42, 2, 393, 394, 7, 9, 2, 2, 394, 396, 5, 84, 43, 2, 395, 393, 3,
	2, 2, 2, 396, 399, 3, 2, 2, 2, 397, 395, 3, 2, 2, 2, 397, 398, 3, 2, 2,
	2, 398, 401, 3, 2, 2, 2, 399, 397, 3, 2, 2, 2, 400, 392, 3, 2, 2, 2, 401,
	404, 3, 2, 2, 2, 402, 400, 3, 2, 2, 2, 402, 403, 3, 2, 2, 2, 403, 406,
	3, 2, 2, 2, 404, 402, 3, 2, 2, 2, 405, 364, 3, 2, 2, 2, 405, 377, 3, 2,
	2, 2, 406, 79, 3, 2, 2, 2, 407, 408, 5, 56, 29, 2, 408, 81, 3, 2, 2, 2,
	409, 410, 7, 11, 2, 2, 410, 411, 5, 92, 47, 2, 411, 412, 7, 12, 2, 2, 412,
	83, 3, 2, 2, 2, 413, 416, 7, 60, 2, 2, 414, 416, 5, 66, 34, 2, 415, 413,
	3, 2, 2, 2, 415, 414, 3, 2, 2, 2, 416, 85, 3, 2, 2, 2, 417, 422, 5, 92,
	47, 2, 418, 419, 7, 10, 2, 2, 419, 421, 5, 92, 47, 2, 420, 418, 3, 2, 2,
	2, 421, 424, 3, 2, 2, 2, 422, 420, 3, 2, 2, 2, 422, 423, 3, 2, 2, 2, 423,
	87, 3, 2, 2, 2, 424, 422, 3, 2, 2, 2, 425, 426, 7, 60, 2, 2, 426, 427,
	5, 90, 46, 2, 427, 89, 3, 2, 2, 2, 428, 437, 7, 13, 2, 2, 429, 434, 5,
	92, 47, 2, 430, 431, 7, 10, 2, 2, 431, 433, 5, 92, 47, 2, 432, 430, 3,
	2, 2, 2, 433, 436, 3, 2, 2, 2, 434, 432, 3, 2, 2, 2, 434, 435, 3, 2, 2,
	2, 435, 438, 3, 2, 2, 2, 436, 434, 3, 2, 2, 2, 437, 429, 3, 2, 2, 2, 437,
	438, 3, 2, 2, 2, 438, 439, 3, 2, 2, 2, 439, 440, 7, 14, 2, 2, 440, 91,
	3, 2, 2, 2, 441, 442, 8, 47, 1, 2, 442, 443, 5, 106, 54, 2, 443, 444, 5,
	92, 47, 22, 444, 462, 3, 2, 2, 2, 445, 462, 5, 88, 45, 2, 446, 447, 7,
	13, 2, 2, 447, 448, 5, 86, 44, 2, 448, 449, 7, 14, 2, 2, 449, 462, 3, 2,
	2, 2, 450, 462, 5, 58, 30, 2, 451, 462, 5, 66, 34, 2, 452, 462, 5, 68,
	35, 2, 453, 462, 5, 70, 36, 2, 454, 462, 5, 64, 33, 2, 455, 462, 5, 60,
	31, 2, 456, 462, 5, 62, 32, 2, 457, 462, 5, 56, 29, 2, 458, 462, 5, 78,
	40, 2, 459, 462, 5, 72, 37, 2, 460, 462, 5, 54, 28, 2, 461, 441, 3, 2,
	2, 2, 461, 445, 3, 2, 2, 2, 461, 446, 3, 2, 2, 2, 461, 450, 3, 2, 2, 2,
	461, 451, 3, 2, 2, 2, 461, 452, 3, 2, 2, 2, 461, 453, 3, 2, 2, 2, 461,
	454, 3, 2, 2, 2, 461, 455, 3, 2, 2, 2, 461, 456, 3, 2, 2, 2, 461, 457,
	3, 2, 2, 2, 461, 458, 3, 2, 2, 2, 461, 459, 3, 2, 2, 2, 461, 460, 3, 2,
	2, 2, 462, 496, 3, 2, 2, 2, 463, 464, 12, 21, 2, 2, 464, 465, 5, 100, 51,
	2, 465, 466, 5, 92, 47, 22, 466, 495, 3, 2, 2, 2, 467, 468, 12, 20, 2,
	2, 468, 469, 5, 102, 52, 2, 469, 470, 5, 92, 47, 21, 470, 495, 3, 2, 2,
	2, 471, 472, 12, 19, 2, 2, 472, 473, 5, 104, 53, 2, 473, 474, 5, 92, 47,
	20, 474, 495, 3, 2, 2, 2, 475, 476, 12, 16, 2, 2, 476, 479, 5, 96, 49,
	2, 477, 480, 5, 98, 50, 2, 478, 480, 5, 100, 51, 2, 479, 477, 3, 2, 2,
	2, 479, 478, 3, 2, 2, 2, 480, 481, 3, 2, 2, 2, 481, 482, 5, 92, 47, 17,
	482, 495, 3, 2, 2, 2, 483, 484, 12, 15, 2, 2, 484, 485, 5, 98, 50, 2, 485,
	486, 5, 92, 47, 16, 486, 495, 3, 2, 2, 2, 487, 488, 12, 14, 2, 2, 488,
	490, 7, 34, 2, 2, 489, 491, 5, 92, 47, 2, 490, 489, 3, 2, 2, 2, 490, 491,
	3, 2, 2, 2, 491, 492, 3, 2, 2, 2, 492, 493, 7, 7, 2, 2, 493, 495, 5, 92,
	47, 15, 494, 463, 3, 2, 2, 2, 494, 467, 3, 2, 2, 2, 494, 471, 3, 2, 2,
	2, 494, 475, 3, 2, 2, 2, 494, 483, 3, 2, 2, 2, 494, 487, 3, 2, 2, 2, 495,
	498, 3, 2, 2, 2, 496, 494, 3, 2, 2, 2, 496, 497, 3, 2, 2, 2, 497, 93, 3,
	2, 2, 2, 498, 496, 3, 2, 2, 2, 499, 500, 5, 92, 47, 2, 500, 502, 7, 34,
	2, 2, 501, 503, 5, 92, 47, 2, 502, 501, 3, 2, 2, 2, 502, 503, 3, 2, 2,
	2, 503, 504, 3, 2, 2, 2, 504, 505, 7, 7, 2, 2, 505, 506, 7, 13, 2, 2, 506,
	507, 5, 12, 7, 2, 507, 508, 7, 14, 2, 2, 508, 528, 3, 2, 2, 2, 509, 510,
	5, 92, 47, 2, 510, 511, 7, 34, 2, 2, 511, 512, 7, 13, 2, 2, 512, 513, 5,
	12, 7, 2, 513, 514, 7, 14, 2, 2, 514, 515, 7, 7, 2, 2, 515, 516, 5, 92,
	47, 2, 516, 528, 3, 2, 2, 2, 517, 518, 5, 92, 47, 2, 518, 519, 7, 34, 2,
	2, 519, 520, 7, 13, 2, 2, 520, 521, 5, 12, 7, 2, 521, 522, 7, 14, 2, 2,
	522, 523, 7, 7, 2, 2, 523, 524, 7, 13, 2, 2, 524, 525, 5, 12, 7, 2, 525,
	526, 7, 14, 2, 2, 526, 528, 3, 2, 2, 2, 527, 499, 3, 2, 2, 2, 527, 509,
	3, 2, 2, 2, 527, 517, 3, 2, 2, 2, 528, 95, 3, 2, 2, 2, 529, 530, 9, 4,
	2, 2, 530, 97, 3, 2, 2, 2, 531, 535, 7, 58, 2, 2, 532, 533, 7, 57, 2, 2,
	533, 535, 7, 58, 2, 2, 534, 531, 3, 2, 2, 2, 534, 532, 3, 2, 2, 2, 535,
	99, 3, 2, 2, 2, 536, 537, 9, 5, 2, 2, 537, 101, 3, 2, 2, 2, 538, 539, 9,
	6, 2, 2, 539, 103, 3, 2, 2, 2, 540, 541, 9, 7, 2, 2, 541, 105, 3, 2, 2,
	2, 542, 543, 9, 8, 2, 2, 543, 107, 3, 2, 2, 2, 53, 113, 120, 124, 128,
	133, 141, 147, 154, 170, 176, 180, 184, 188, 197, 201, 209, 214, 234, 245,
	254, 267, 269, 291, 301, 307, 311, 321, 324, 327, 345, 350, 362, 370, 375,
	384, 389, 397, 402, 405, 415, 422, 434, 437, 461, 479, 490, 494, 496, 502,
	527, 534,
}
var deserializer = antlr.NewATNDeserializer(nil)
var deserializedATN = deserializer.DeserializeFromUInt16(parserATN)

var literalNames = []string{
	"", "", "", "", "", "':'", "';'", "'.'", "','", "'['", "']'", "'('", "')'",
	"'{'", "'}'", "'>'", "'<'", "'=='", "'>='", "'<='", "'!='", "'+'", "'-'",
	"'--'", "'++'", "'*'", "'/'", "'%'", "", "", "", "'='", "'?'", "'!~'",
	"'=~'", "'FOR'", "'RETURN'", "'DISTINCT'", "'FILTER'", "'SORT'", "'LIMIT'",
	"'LET'", "'COLLECT'", "", "'NONE'", "'NULL'", "", "'INTO'", "'KEEP'", "'WITH'",
	"'COUNT'", "'ALL'", "'ANY'", "'AGGREGATE'", "'LIKE'", "", "'IN'", "'@'",
}
var symbolicNames = []string{
	"", "MultiLineComment", "SingleLineComment", "WhiteSpaces", "LineTerminator",
	"Colon", "SemiColon", "Dot", "Comma", "OpenBracket", "CloseBracket", "OpenParen",
	"CloseParen", "OpenBrace", "CloseBrace", "Gt", "Lt", "Eq", "Gte", "Lte",
	"Neq", "Plus", "Minus", "MinusMinus", "PlusPlus", "Multi", "Div", "Mod",
	"And", "Or", "Range", "Assign", "QuestionMark", "RegexNotMatch", "RegexMatch",
	"For", "Return", "Distinct", "Filter", "Sort", "Limit", "Let", "Collect",
	"SortDirection", "None", "Null", "BooleanLiteral", "Into", "Keep", "With",
	"Count", "All", "Any", "Aggregate", "Like", "Not", "In", "Param", "Identifier",
	"StringLiteral", "TemplateStringLiteral", "IntegerLiteral", "FloatLiteral",
}

var ruleNames = []string{
	"program", "body", "bodyStatement", "bodyExpression", "returnExpression",
	"forExpression", "forExpressionValueVariable", "forExpressionKeyVariable",
	"forExpressionSource", "forExpressionClause", "forExpressionStatement",
	"forExpressionBody", "forExpressionReturn", "filterClause", "limitClause",
	"limitClauseValue", "sortClause", "sortClauseExpression", "collectClause",
	"collectSelector", "collectGrouping", "collectAggregator", "collectAggregateSelector",
	"collectGroupVariable", "collectCounter", "variableDeclaration", "param",
	"variable", "rangeOperator", "arrayLiteral", "objectLiteral", "booleanLiteral",
	"stringLiteral", "integerLiteral", "floatLiteral", "noneLiteral", "arrayElementList",
	"propertyAssignment", "memberExpression", "shorthandPropertyName", "computedPropertyName",
	"propertyName", "expressionSequence", "functionCallExpression", "arguments",
	"expression", "forTernaryExpression", "arrayOperator", "inOperator", "equalityOperator",
	"logicalOperator", "mathOperator", "unaryOperator",
}
var decisionToDFA = make([]*antlr.DFA, len(deserializedATN.DecisionToState))

func init() {
	for index, ds := range deserializedATN.DecisionToState {
		decisionToDFA[index] = antlr.NewDFA(ds, index)
	}
}

type FqlParser struct {
	*antlr.BaseParser
}

func NewFqlParser(input antlr.TokenStream) *FqlParser {
	this := new(FqlParser)

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
	FqlParserEOF                   = antlr.TokenEOF
	FqlParserMultiLineComment      = 1
	FqlParserSingleLineComment     = 2
	FqlParserWhiteSpaces           = 3
	FqlParserLineTerminator        = 4
	FqlParserColon                 = 5
	FqlParserSemiColon             = 6
	FqlParserDot                   = 7
	FqlParserComma                 = 8
	FqlParserOpenBracket           = 9
	FqlParserCloseBracket          = 10
	FqlParserOpenParen             = 11
	FqlParserCloseParen            = 12
	FqlParserOpenBrace             = 13
	FqlParserCloseBrace            = 14
	FqlParserGt                    = 15
	FqlParserLt                    = 16
	FqlParserEq                    = 17
	FqlParserGte                   = 18
	FqlParserLte                   = 19
	FqlParserNeq                   = 20
	FqlParserPlus                  = 21
	FqlParserMinus                 = 22
	FqlParserMinusMinus            = 23
	FqlParserPlusPlus              = 24
	FqlParserMulti                 = 25
	FqlParserDiv                   = 26
	FqlParserMod                   = 27
	FqlParserAnd                   = 28
	FqlParserOr                    = 29
	FqlParserRange                 = 30
	FqlParserAssign                = 31
	FqlParserQuestionMark          = 32
	FqlParserRegexNotMatch         = 33
	FqlParserRegexMatch            = 34
	FqlParserFor                   = 35
	FqlParserReturn                = 36
	FqlParserDistinct              = 37
	FqlParserFilter                = 38
	FqlParserSort                  = 39
	FqlParserLimit                 = 40
	FqlParserLet                   = 41
	FqlParserCollect               = 42
	FqlParserSortDirection         = 43
	FqlParserNone                  = 44
	FqlParserNull                  = 45
	FqlParserBooleanLiteral        = 46
	FqlParserInto                  = 47
	FqlParserKeep                  = 48
	FqlParserWith                  = 49
	FqlParserCount                 = 50
	FqlParserAll                   = 51
	FqlParserAny                   = 52
	FqlParserAggregate             = 53
	FqlParserLike                  = 54
	FqlParserNot                   = 55
	FqlParserIn                    = 56
	FqlParserParam                 = 57
	FqlParserIdentifier            = 58
	FqlParserStringLiteral         = 59
	FqlParserTemplateStringLiteral = 60
	FqlParserIntegerLiteral        = 61
	FqlParserFloatLiteral          = 62
)

// FqlParser rules.
const (
	FqlParserRULE_program                    = 0
	FqlParserRULE_body                       = 1
	FqlParserRULE_bodyStatement              = 2
	FqlParserRULE_bodyExpression             = 3
	FqlParserRULE_returnExpression           = 4
	FqlParserRULE_forExpression              = 5
	FqlParserRULE_forExpressionValueVariable = 6
	FqlParserRULE_forExpressionKeyVariable   = 7
	FqlParserRULE_forExpressionSource        = 8
	FqlParserRULE_forExpressionClause        = 9
	FqlParserRULE_forExpressionStatement     = 10
	FqlParserRULE_forExpressionBody          = 11
	FqlParserRULE_forExpressionReturn        = 12
	FqlParserRULE_filterClause               = 13
	FqlParserRULE_limitClause                = 14
	FqlParserRULE_limitClauseValue           = 15
	FqlParserRULE_sortClause                 = 16
	FqlParserRULE_sortClauseExpression       = 17
	FqlParserRULE_collectClause              = 18
	FqlParserRULE_collectSelector            = 19
	FqlParserRULE_collectGrouping            = 20
	FqlParserRULE_collectAggregator          = 21
	FqlParserRULE_collectAggregateSelector   = 22
	FqlParserRULE_collectGroupVariable       = 23
	FqlParserRULE_collectCounter             = 24
	FqlParserRULE_variableDeclaration        = 25
	FqlParserRULE_param                      = 26
	FqlParserRULE_variable                   = 27
	FqlParserRULE_rangeOperator              = 28
	FqlParserRULE_arrayLiteral               = 29
	FqlParserRULE_objectLiteral              = 30
	FqlParserRULE_booleanLiteral             = 31
	FqlParserRULE_stringLiteral              = 32
	FqlParserRULE_integerLiteral             = 33
	FqlParserRULE_floatLiteral               = 34
	FqlParserRULE_noneLiteral                = 35
	FqlParserRULE_arrayElementList           = 36
	FqlParserRULE_propertyAssignment         = 37
	FqlParserRULE_memberExpression           = 38
	FqlParserRULE_shorthandPropertyName      = 39
	FqlParserRULE_computedPropertyName       = 40
	FqlParserRULE_propertyName               = 41
	FqlParserRULE_expressionSequence         = 42
	FqlParserRULE_functionCallExpression     = 43
	FqlParserRULE_arguments                  = 44
	FqlParserRULE_expression                 = 45
	FqlParserRULE_forTernaryExpression       = 46
	FqlParserRULE_arrayOperator              = 47
	FqlParserRULE_inOperator                 = 48
	FqlParserRULE_equalityOperator           = 49
	FqlParserRULE_logicalOperator            = 50
	FqlParserRULE_mathOperator               = 51
	FqlParserRULE_unaryOperator              = 52
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

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(106)
		p.Body()
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
	p.EnterRule(localctx, 2, FqlParserRULE_body)
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
	p.SetState(111)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == FqlParserLet || _la == FqlParserIdentifier {
		{
			p.SetState(108)
			p.BodyStatement()
		}

		p.SetState(113)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(114)
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
	p.EnterRule(localctx, 4, FqlParserRULE_bodyStatement)

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

	p.SetState(118)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case FqlParserIdentifier:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(116)
			p.FunctionCallExpression()
		}

	case FqlParserLet:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(117)
			p.VariableDeclaration()
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
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
	p.EnterRule(localctx, 6, FqlParserRULE_bodyExpression)

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

	p.SetState(122)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case FqlParserReturn:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(120)
			p.ReturnExpression()
		}

	case FqlParserFor:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(121)
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
	p.EnterRule(localctx, 8, FqlParserRULE_returnExpression)
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

	p.SetState(139)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 5, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(124)
			p.Match(FqlParserReturn)
		}
		p.SetState(126)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		if _la == FqlParserDistinct {
			{
				p.SetState(125)
				p.Match(FqlParserDistinct)
			}

		}
		{
			p.SetState(128)
			p.expression(0)
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(129)
			p.Match(FqlParserReturn)
		}
		p.SetState(131)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		if _la == FqlParserDistinct {
			{
				p.SetState(130)
				p.Match(FqlParserDistinct)
			}

		}
		{
			p.SetState(133)
			p.Match(FqlParserOpenParen)
		}
		{
			p.SetState(134)
			p.ForExpression()
		}
		{
			p.SetState(135)
			p.Match(FqlParserCloseParen)
		}

	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(137)
			p.Match(FqlParserReturn)
		}
		{
			p.SetState(138)
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
	p.EnterRule(localctx, 10, FqlParserRULE_forExpression)
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
		p.SetState(141)
		p.Match(FqlParserFor)
	}
	{
		p.SetState(142)
		p.ForExpressionValueVariable()
	}
	p.SetState(145)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == FqlParserComma {
		{
			p.SetState(143)
			p.Match(FqlParserComma)
		}
		{
			p.SetState(144)
			p.ForExpressionKeyVariable()
		}

	}
	{
		p.SetState(147)
		p.Match(FqlParserIn)
	}
	{
		p.SetState(148)
		p.ForExpressionSource()
	}
	p.SetState(152)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for ((_la-38)&-(0x1f+1)) == 0 && ((1<<uint((_la-38)))&((1<<(FqlParserFilter-38))|(1<<(FqlParserSort-38))|(1<<(FqlParserLimit-38))|(1<<(FqlParserLet-38))|(1<<(FqlParserCollect-38))|(1<<(FqlParserIdentifier-38)))) != 0 {
		{
			p.SetState(149)
			p.ForExpressionBody()
		}

		p.SetState(154)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(155)
		p.ForExpressionReturn()
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
	p.EnterRule(localctx, 12, FqlParserRULE_forExpressionValueVariable)

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
		p.SetState(157)
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
	p.EnterRule(localctx, 14, FqlParserRULE_forExpressionKeyVariable)

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
		p.SetState(159)
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
	p.EnterRule(localctx, 16, FqlParserRULE_forExpressionSource)

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
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 8, p.GetParserRuleContext()) {
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
			p.ArrayLiteral()
		}

	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(163)
			p.ObjectLiteral()
		}

	case 4:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(164)
			p.Variable()
		}

	case 5:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(165)
			p.MemberExpression()
		}

	case 6:
		p.EnterOuterAlt(localctx, 6)
		{
			p.SetState(166)
			p.RangeOperator()
		}

	case 7:
		p.EnterOuterAlt(localctx, 7)
		{
			p.SetState(167)
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
	p.EnterRule(localctx, 18, FqlParserRULE_forExpressionClause)

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

	p.SetState(174)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case FqlParserLimit:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(170)
			p.LimitClause()
		}

	case FqlParserSort:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(171)
			p.SortClause()
		}

	case FqlParserFilter:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(172)
			p.FilterClause()
		}

	case FqlParserCollect:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(173)
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
	p.EnterRule(localctx, 20, FqlParserRULE_forExpressionStatement)

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

	p.SetState(178)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case FqlParserLet:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(176)
			p.VariableDeclaration()
		}

	case FqlParserIdentifier:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(177)
			p.FunctionCallExpression()
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
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
	p.EnterRule(localctx, 22, FqlParserRULE_forExpressionBody)

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

	p.SetState(182)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case FqlParserLet, FqlParserIdentifier:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(180)
			p.ForExpressionStatement()
		}

	case FqlParserFilter, FqlParserSort, FqlParserLimit, FqlParserCollect:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(181)
			p.ForExpressionClause()
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
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
	p.EnterRule(localctx, 24, FqlParserRULE_forExpressionReturn)

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

	p.SetState(186)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case FqlParserReturn:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(184)
			p.ReturnExpression()
		}

	case FqlParserFor:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(185)
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
	p.EnterRule(localctx, 26, FqlParserRULE_filterClause)

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
		p.SetState(188)
		p.Match(FqlParserFilter)
	}
	{
		p.SetState(189)
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
	p.EnterRule(localctx, 28, FqlParserRULE_limitClause)
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
		p.SetState(191)
		p.Match(FqlParserLimit)
	}
	{
		p.SetState(192)
		p.LimitClauseValue()
	}
	p.SetState(195)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == FqlParserComma {
		{
			p.SetState(193)
			p.Match(FqlParserComma)
		}
		{
			p.SetState(194)
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
	p.EnterRule(localctx, 30, FqlParserRULE_limitClauseValue)

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

	p.SetState(199)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case FqlParserIntegerLiteral:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(197)
			p.Match(FqlParserIntegerLiteral)
		}

	case FqlParserParam:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(198)
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
	p.EnterRule(localctx, 32, FqlParserRULE_sortClause)
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
		p.SetState(201)
		p.Match(FqlParserSort)
	}
	{
		p.SetState(202)
		p.SortClauseExpression()
	}
	p.SetState(207)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == FqlParserComma {
		{
			p.SetState(203)
			p.Match(FqlParserComma)
		}
		{
			p.SetState(204)
			p.SortClauseExpression()
		}

		p.SetState(209)
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
	p.EnterRule(localctx, 34, FqlParserRULE_sortClauseExpression)
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
		p.SetState(210)
		p.expression(0)
	}
	p.SetState(212)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == FqlParserSortDirection {
		{
			p.SetState(211)
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
	p.EnterRule(localctx, 36, FqlParserRULE_collectClause)

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

	p.SetState(232)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 17, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(214)
			p.Match(FqlParserCollect)
		}
		{
			p.SetState(215)
			p.CollectCounter()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(216)
			p.Match(FqlParserCollect)
		}
		{
			p.SetState(217)
			p.CollectAggregator()
		}

	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(218)
			p.Match(FqlParserCollect)
		}
		{
			p.SetState(219)
			p.CollectGrouping()
		}
		{
			p.SetState(220)
			p.CollectAggregator()
		}

	case 4:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(222)
			p.Match(FqlParserCollect)
		}
		{
			p.SetState(223)
			p.CollectGrouping()
		}
		{
			p.SetState(224)
			p.CollectGroupVariable()
		}

	case 5:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(226)
			p.Match(FqlParserCollect)
		}
		{
			p.SetState(227)
			p.CollectGrouping()
		}
		{
			p.SetState(228)
			p.CollectCounter()
		}

	case 6:
		p.EnterOuterAlt(localctx, 6)
		{
			p.SetState(230)
			p.Match(FqlParserCollect)
		}
		{
			p.SetState(231)
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
	p.EnterRule(localctx, 38, FqlParserRULE_collectSelector)

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
		p.SetState(234)
		p.Match(FqlParserIdentifier)
	}
	{
		p.SetState(235)
		p.Match(FqlParserAssign)
	}
	{
		p.SetState(236)
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
	p.EnterRule(localctx, 40, FqlParserRULE_collectGrouping)
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
		p.SetState(238)
		p.CollectSelector()
	}
	p.SetState(243)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == FqlParserComma {
		{
			p.SetState(239)
			p.Match(FqlParserComma)
		}
		{
			p.SetState(240)
			p.CollectSelector()
		}

		p.SetState(245)
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
	p.EnterRule(localctx, 42, FqlParserRULE_collectAggregator)
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
		p.SetState(246)
		p.Match(FqlParserAggregate)
	}
	{
		p.SetState(247)
		p.CollectAggregateSelector()
	}
	p.SetState(252)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == FqlParserComma {
		{
			p.SetState(248)
			p.Match(FqlParserComma)
		}
		{
			p.SetState(249)
			p.CollectAggregateSelector()
		}

		p.SetState(254)
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
	p.EnterRule(localctx, 44, FqlParserRULE_collectAggregateSelector)

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
		p.SetState(255)
		p.Match(FqlParserIdentifier)
	}
	{
		p.SetState(256)
		p.Match(FqlParserAssign)
	}
	{
		p.SetState(257)
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
	p.EnterRule(localctx, 46, FqlParserRULE_collectGroupVariable)
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

	p.SetState(267)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 21, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(259)
			p.Match(FqlParserInto)
		}
		{
			p.SetState(260)
			p.CollectSelector()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(261)
			p.Match(FqlParserInto)
		}
		{
			p.SetState(262)
			p.Match(FqlParserIdentifier)
		}
		p.SetState(265)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		if _la == FqlParserKeep {
			{
				p.SetState(263)
				p.Match(FqlParserKeep)
			}
			{
				p.SetState(264)
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
	p.EnterRule(localctx, 48, FqlParserRULE_collectCounter)

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
		p.SetState(269)
		p.Match(FqlParserWith)
	}
	{
		p.SetState(270)
		p.Match(FqlParserCount)
	}
	{
		p.SetState(271)
		p.Match(FqlParserInto)
	}
	{
		p.SetState(272)
		p.Match(FqlParserIdentifier)
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
	p.EnterRule(localctx, 50, FqlParserRULE_variableDeclaration)

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

	p.SetState(289)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 22, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(274)
			p.Match(FqlParserLet)
		}
		{
			p.SetState(275)
			p.Match(FqlParserIdentifier)
		}
		{
			p.SetState(276)
			p.Match(FqlParserAssign)
		}
		{
			p.SetState(277)
			p.expression(0)
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(278)
			p.Match(FqlParserLet)
		}
		{
			p.SetState(279)
			p.Match(FqlParserIdentifier)
		}
		{
			p.SetState(280)
			p.Match(FqlParserAssign)
		}
		{
			p.SetState(281)
			p.Match(FqlParserOpenParen)
		}
		{
			p.SetState(282)
			p.ForExpression()
		}
		{
			p.SetState(283)
			p.Match(FqlParserCloseParen)
		}

	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(285)
			p.Match(FqlParserLet)
		}
		{
			p.SetState(286)
			p.Match(FqlParserIdentifier)
		}
		{
			p.SetState(287)
			p.Match(FqlParserAssign)
		}
		{
			p.SetState(288)
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
	p.EnterRule(localctx, 52, FqlParserRULE_param)

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
		p.SetState(291)
		p.Match(FqlParserParam)
	}
	{
		p.SetState(292)
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
	p.EnterRule(localctx, 54, FqlParserRULE_variable)

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
		p.SetState(294)
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
	p.EnterRule(localctx, 56, FqlParserRULE_rangeOperator)

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
	p.SetState(299)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case FqlParserIntegerLiteral:
		{
			p.SetState(296)
			p.IntegerLiteral()
		}

	case FqlParserIdentifier:
		{
			p.SetState(297)
			p.Variable()
		}

	case FqlParserParam:
		{
			p.SetState(298)
			p.Param()
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}
	{
		p.SetState(301)
		p.Match(FqlParserRange)
	}
	p.SetState(305)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case FqlParserIntegerLiteral:
		{
			p.SetState(302)
			p.IntegerLiteral()
		}

	case FqlParserIdentifier:
		{
			p.SetState(303)
			p.Variable()
		}

	case FqlParserParam:
		{
			p.SetState(304)
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
	p.EnterRule(localctx, 58, FqlParserRULE_arrayLiteral)
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
		p.SetState(307)
		p.Match(FqlParserOpenBracket)
	}
	p.SetState(309)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if (((_la)&-(0x1f+1)) == 0 && ((1<<uint(_la))&((1<<FqlParserOpenBracket)|(1<<FqlParserOpenParen)|(1<<FqlParserOpenBrace)|(1<<FqlParserPlus)|(1<<FqlParserMinus))) != 0) || (((_la-44)&-(0x1f+1)) == 0 && ((1<<uint((_la-44)))&((1<<(FqlParserNone-44))|(1<<(FqlParserNull-44))|(1<<(FqlParserBooleanLiteral-44))|(1<<(FqlParserLike-44))|(1<<(FqlParserNot-44))|(1<<(FqlParserParam-44))|(1<<(FqlParserIdentifier-44))|(1<<(FqlParserStringLiteral-44))|(1<<(FqlParserTemplateStringLiteral-44))|(1<<(FqlParserIntegerLiteral-44))|(1<<(FqlParserFloatLiteral-44)))) != 0) {
		{
			p.SetState(308)
			p.ArrayElementList()
		}

	}
	{
		p.SetState(311)
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
	p.EnterRule(localctx, 60, FqlParserRULE_objectLiteral)
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
		p.SetState(313)
		p.Match(FqlParserOpenBrace)
	}
	p.SetState(322)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == FqlParserOpenBracket || (((_la-58)&-(0x1f+1)) == 0 && ((1<<uint((_la-58)))&((1<<(FqlParserIdentifier-58))|(1<<(FqlParserStringLiteral-58))|(1<<(FqlParserTemplateStringLiteral-58)))) != 0) {
		{
			p.SetState(314)
			p.PropertyAssignment()
		}
		p.SetState(319)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 26, p.GetParserRuleContext())

		for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
			if _alt == 1 {
				{
					p.SetState(315)
					p.Match(FqlParserComma)
				}
				{
					p.SetState(316)
					p.PropertyAssignment()
				}

			}
			p.SetState(321)
			p.GetErrorHandler().Sync(p)
			_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 26, p.GetParserRuleContext())
		}

	}
	p.SetState(325)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == FqlParserComma {
		{
			p.SetState(324)
			p.Match(FqlParserComma)
		}

	}
	{
		p.SetState(327)
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
	p.EnterRule(localctx, 62, FqlParserRULE_booleanLiteral)

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
		p.SetState(329)
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

func (s *StringLiteralContext) TemplateStringLiteral() antlr.TerminalNode {
	return s.GetToken(FqlParserTemplateStringLiteral, 0)
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
	p.EnterRule(localctx, 64, FqlParserRULE_stringLiteral)
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
		p.SetState(331)
		_la = p.GetTokenStream().LA(1)

		if !(_la == FqlParserStringLiteral || _la == FqlParserTemplateStringLiteral) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
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
	p.EnterRule(localctx, 66, FqlParserRULE_integerLiteral)

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
		p.SetState(333)
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
	p.EnterRule(localctx, 68, FqlParserRULE_floatLiteral)

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
		p.SetState(335)
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
	p.EnterRule(localctx, 70, FqlParserRULE_noneLiteral)
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
		p.SetState(337)
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
	p.EnterRule(localctx, 72, FqlParserRULE_arrayElementList)
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
		p.SetState(339)
		p.expression(0)
	}
	p.SetState(348)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == FqlParserComma {
		p.SetState(341)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		for ok := true; ok; ok = _la == FqlParserComma {
			{
				p.SetState(340)
				p.Match(FqlParserComma)
			}

			p.SetState(343)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)
		}
		{
			p.SetState(345)
			p.expression(0)
		}

		p.SetState(350)
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
	p.EnterRule(localctx, 74, FqlParserRULE_propertyAssignment)

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

	p.SetState(360)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 31, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(351)
			p.PropertyName()
		}
		{
			p.SetState(352)
			p.Match(FqlParserColon)
		}
		{
			p.SetState(353)
			p.expression(0)
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(355)
			p.ComputedPropertyName()
		}
		{
			p.SetState(356)
			p.Match(FqlParserColon)
		}
		{
			p.SetState(357)
			p.expression(0)
		}

	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(359)
			p.ShorthandPropertyName()
		}

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

func (s *MemberExpressionContext) Identifier() antlr.TerminalNode {
	return s.GetToken(FqlParserIdentifier, 0)
}

func (s *MemberExpressionContext) AllDot() []antlr.TerminalNode {
	return s.GetTokens(FqlParserDot)
}

func (s *MemberExpressionContext) Dot(i int) antlr.TerminalNode {
	return s.GetToken(FqlParserDot, i)
}

func (s *MemberExpressionContext) AllPropertyName() []IPropertyNameContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IPropertyNameContext)(nil)).Elem())
	var tst = make([]IPropertyNameContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IPropertyNameContext)
		}
	}

	return tst
}

func (s *MemberExpressionContext) PropertyName(i int) IPropertyNameContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IPropertyNameContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IPropertyNameContext)
}

func (s *MemberExpressionContext) AllComputedPropertyName() []IComputedPropertyNameContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IComputedPropertyNameContext)(nil)).Elem())
	var tst = make([]IComputedPropertyNameContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IComputedPropertyNameContext)
		}
	}

	return tst
}

func (s *MemberExpressionContext) ComputedPropertyName(i int) IComputedPropertyNameContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IComputedPropertyNameContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IComputedPropertyNameContext)
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
	p.EnterRule(localctx, 76, FqlParserRULE_memberExpression)

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

	p.SetState(403)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 38, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(362)
			p.Match(FqlParserIdentifier)
		}
		p.SetState(371)
		p.GetErrorHandler().Sync(p)
		_alt = 1
		for ok := true; ok; ok = _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
			switch _alt {
			case 1:
				{
					p.SetState(363)
					p.Match(FqlParserDot)
				}
				{
					p.SetState(364)
					p.PropertyName()
				}
				p.SetState(368)
				p.GetErrorHandler().Sync(p)
				_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 32, p.GetParserRuleContext())

				for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
					if _alt == 1 {
						{
							p.SetState(365)
							p.ComputedPropertyName()
						}

					}
					p.SetState(370)
					p.GetErrorHandler().Sync(p)
					_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 32, p.GetParserRuleContext())
				}

			default:
				panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
			}

			p.SetState(373)
			p.GetErrorHandler().Sync(p)
			_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 33, p.GetParserRuleContext())
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(375)
			p.Match(FqlParserIdentifier)
		}
		{
			p.SetState(376)
			p.ComputedPropertyName()
		}
		p.SetState(387)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 35, p.GetParserRuleContext())

		for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
			if _alt == 1 {
				{
					p.SetState(377)
					p.Match(FqlParserDot)
				}
				{
					p.SetState(378)
					p.PropertyName()
				}
				p.SetState(382)
				p.GetErrorHandler().Sync(p)
				_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 34, p.GetParserRuleContext())

				for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
					if _alt == 1 {
						{
							p.SetState(379)
							p.ComputedPropertyName()
						}

					}
					p.SetState(384)
					p.GetErrorHandler().Sync(p)
					_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 34, p.GetParserRuleContext())
				}

			}
			p.SetState(389)
			p.GetErrorHandler().Sync(p)
			_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 35, p.GetParserRuleContext())
		}
		p.SetState(400)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 37, p.GetParserRuleContext())

		for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
			if _alt == 1 {
				{
					p.SetState(390)
					p.ComputedPropertyName()
				}
				p.SetState(395)
				p.GetErrorHandler().Sync(p)
				_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 36, p.GetParserRuleContext())

				for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
					if _alt == 1 {
						{
							p.SetState(391)
							p.Match(FqlParserDot)
						}
						{
							p.SetState(392)
							p.PropertyName()
						}

					}
					p.SetState(397)
					p.GetErrorHandler().Sync(p)
					_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 36, p.GetParserRuleContext())
				}

			}
			p.SetState(402)
			p.GetErrorHandler().Sync(p)
			_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 37, p.GetParserRuleContext())
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
	p.EnterRule(localctx, 78, FqlParserRULE_shorthandPropertyName)

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
		p.SetState(405)
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
	p.EnterRule(localctx, 80, FqlParserRULE_computedPropertyName)

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
		p.SetState(407)
		p.Match(FqlParserOpenBracket)
	}
	{
		p.SetState(408)
		p.expression(0)
	}
	{
		p.SetState(409)
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
	p.EnterRule(localctx, 82, FqlParserRULE_propertyName)

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

	p.SetState(413)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case FqlParserIdentifier:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(411)
			p.Match(FqlParserIdentifier)
		}

	case FqlParserStringLiteral, FqlParserTemplateStringLiteral:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(412)
			p.StringLiteral()
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}

	return localctx
}

// IExpressionSequenceContext is an interface to support dynamic dispatch.
type IExpressionSequenceContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsExpressionSequenceContext differentiates from other interfaces.
	IsExpressionSequenceContext()
}

type ExpressionSequenceContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyExpressionSequenceContext() *ExpressionSequenceContext {
	var p = new(ExpressionSequenceContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = FqlParserRULE_expressionSequence
	return p
}

func (*ExpressionSequenceContext) IsExpressionSequenceContext() {}

func NewExpressionSequenceContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ExpressionSequenceContext {
	var p = new(ExpressionSequenceContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_expressionSequence

	return p
}

func (s *ExpressionSequenceContext) GetParser() antlr.Parser { return s.parser }

func (s *ExpressionSequenceContext) AllExpression() []IExpressionContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IExpressionContext)(nil)).Elem())
	var tst = make([]IExpressionContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IExpressionContext)
		}
	}

	return tst
}

func (s *ExpressionSequenceContext) Expression(i int) IExpressionContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExpressionContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *ExpressionSequenceContext) AllComma() []antlr.TerminalNode {
	return s.GetTokens(FqlParserComma)
}

func (s *ExpressionSequenceContext) Comma(i int) antlr.TerminalNode {
	return s.GetToken(FqlParserComma, i)
}

func (s *ExpressionSequenceContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ExpressionSequenceContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ExpressionSequenceContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterExpressionSequence(s)
	}
}

func (s *ExpressionSequenceContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitExpressionSequence(s)
	}
}

func (s *ExpressionSequenceContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitExpressionSequence(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) ExpressionSequence() (localctx IExpressionSequenceContext) {
	localctx = NewExpressionSequenceContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 84, FqlParserRULE_expressionSequence)
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
		p.SetState(415)
		p.expression(0)
	}
	p.SetState(420)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == FqlParserComma {
		{
			p.SetState(416)
			p.Match(FqlParserComma)
		}
		{
			p.SetState(417)
			p.expression(0)
		}

		p.SetState(422)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
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

func (s *FunctionCallExpressionContext) Identifier() antlr.TerminalNode {
	return s.GetToken(FqlParserIdentifier, 0)
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
	p.EnterRule(localctx, 86, FqlParserRULE_functionCallExpression)

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
		p.SetState(423)
		p.Match(FqlParserIdentifier)
	}
	{
		p.SetState(424)
		p.Arguments()
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
	p.EnterRule(localctx, 88, FqlParserRULE_arguments)
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
		p.SetState(426)
		p.Match(FqlParserOpenParen)
	}
	p.SetState(435)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if (((_la)&-(0x1f+1)) == 0 && ((1<<uint(_la))&((1<<FqlParserOpenBracket)|(1<<FqlParserOpenParen)|(1<<FqlParserOpenBrace)|(1<<FqlParserPlus)|(1<<FqlParserMinus))) != 0) || (((_la-44)&-(0x1f+1)) == 0 && ((1<<uint((_la-44)))&((1<<(FqlParserNone-44))|(1<<(FqlParserNull-44))|(1<<(FqlParserBooleanLiteral-44))|(1<<(FqlParserLike-44))|(1<<(FqlParserNot-44))|(1<<(FqlParserParam-44))|(1<<(FqlParserIdentifier-44))|(1<<(FqlParserStringLiteral-44))|(1<<(FqlParserTemplateStringLiteral-44))|(1<<(FqlParserIntegerLiteral-44))|(1<<(FqlParserFloatLiteral-44)))) != 0) {
		{
			p.SetState(427)
			p.expression(0)
		}
		p.SetState(432)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		for _la == FqlParserComma {
			{
				p.SetState(428)
				p.Match(FqlParserComma)
			}
			{
				p.SetState(429)
				p.expression(0)
			}

			p.SetState(434)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)
		}

	}
	{
		p.SetState(437)
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

func (s *ExpressionContext) OpenParen() antlr.TerminalNode {
	return s.GetToken(FqlParserOpenParen, 0)
}

func (s *ExpressionContext) ExpressionSequence() IExpressionSequenceContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExpressionSequenceContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IExpressionSequenceContext)
}

func (s *ExpressionContext) CloseParen() antlr.TerminalNode {
	return s.GetToken(FqlParserCloseParen, 0)
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

func (s *ExpressionContext) EqualityOperator() IEqualityOperatorContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IEqualityOperatorContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IEqualityOperatorContext)
}

func (s *ExpressionContext) LogicalOperator() ILogicalOperatorContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ILogicalOperatorContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ILogicalOperatorContext)
}

func (s *ExpressionContext) MathOperator() IMathOperatorContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IMathOperatorContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IMathOperatorContext)
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
	_startState := 90
	p.EnterRecursionRule(localctx, 90, FqlParserRULE_expression, _p)
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
	p.SetState(459)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 43, p.GetParserRuleContext()) {
	case 1:
		{
			p.SetState(440)
			p.UnaryOperator()
		}
		{
			p.SetState(441)
			p.expression(20)
		}

	case 2:
		{
			p.SetState(443)
			p.FunctionCallExpression()
		}

	case 3:
		{
			p.SetState(444)
			p.Match(FqlParserOpenParen)
		}
		{
			p.SetState(445)
			p.ExpressionSequence()
		}
		{
			p.SetState(446)
			p.Match(FqlParserCloseParen)
		}

	case 4:
		{
			p.SetState(448)
			p.RangeOperator()
		}

	case 5:
		{
			p.SetState(449)
			p.StringLiteral()
		}

	case 6:
		{
			p.SetState(450)
			p.IntegerLiteral()
		}

	case 7:
		{
			p.SetState(451)
			p.FloatLiteral()
		}

	case 8:
		{
			p.SetState(452)
			p.BooleanLiteral()
		}

	case 9:
		{
			p.SetState(453)
			p.ArrayLiteral()
		}

	case 10:
		{
			p.SetState(454)
			p.ObjectLiteral()
		}

	case 11:
		{
			p.SetState(455)
			p.Variable()
		}

	case 12:
		{
			p.SetState(456)
			p.MemberExpression()
		}

	case 13:
		{
			p.SetState(457)
			p.NoneLiteral()
		}

	case 14:
		{
			p.SetState(458)
			p.Param()
		}

	}
	p.GetParserRuleContext().SetStop(p.GetTokenStream().LT(-1))
	p.SetState(494)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 47, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			if p.GetParseListeners() != nil {
				p.TriggerExitRuleEvent()
			}
			_prevctx = localctx
			p.SetState(492)
			p.GetErrorHandler().Sync(p)
			switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 46, p.GetParserRuleContext()) {
			case 1:
				localctx = NewExpressionContext(p, _parentctx, _parentState)
				p.PushNewRecursionContext(localctx, _startState, FqlParserRULE_expression)
				p.SetState(461)

				if !(p.Precpred(p.GetParserRuleContext(), 19)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 19)", ""))
				}
				{
					p.SetState(462)
					p.EqualityOperator()
				}
				{
					p.SetState(463)
					p.expression(20)
				}

			case 2:
				localctx = NewExpressionContext(p, _parentctx, _parentState)
				p.PushNewRecursionContext(localctx, _startState, FqlParserRULE_expression)
				p.SetState(465)

				if !(p.Precpred(p.GetParserRuleContext(), 18)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 18)", ""))
				}
				{
					p.SetState(466)
					p.LogicalOperator()
				}
				{
					p.SetState(467)
					p.expression(19)
				}

			case 3:
				localctx = NewExpressionContext(p, _parentctx, _parentState)
				p.PushNewRecursionContext(localctx, _startState, FqlParserRULE_expression)
				p.SetState(469)

				if !(p.Precpred(p.GetParserRuleContext(), 17)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 17)", ""))
				}
				{
					p.SetState(470)
					p.MathOperator()
				}
				{
					p.SetState(471)
					p.expression(18)
				}

			case 4:
				localctx = NewExpressionContext(p, _parentctx, _parentState)
				p.PushNewRecursionContext(localctx, _startState, FqlParserRULE_expression)
				p.SetState(473)

				if !(p.Precpred(p.GetParserRuleContext(), 14)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 14)", ""))
				}
				{
					p.SetState(474)
					p.ArrayOperator()
				}
				p.SetState(477)
				p.GetErrorHandler().Sync(p)

				switch p.GetTokenStream().LA(1) {
				case FqlParserNot, FqlParserIn:
					{
						p.SetState(475)
						p.InOperator()
					}

				case FqlParserGt, FqlParserLt, FqlParserEq, FqlParserGte, FqlParserLte, FqlParserNeq:
					{
						p.SetState(476)
						p.EqualityOperator()
					}

				default:
					panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
				}
				{
					p.SetState(479)
					p.expression(15)
				}

			case 5:
				localctx = NewExpressionContext(p, _parentctx, _parentState)
				p.PushNewRecursionContext(localctx, _startState, FqlParserRULE_expression)
				p.SetState(481)

				if !(p.Precpred(p.GetParserRuleContext(), 13)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 13)", ""))
				}
				{
					p.SetState(482)
					p.InOperator()
				}
				{
					p.SetState(483)
					p.expression(14)
				}

			case 6:
				localctx = NewExpressionContext(p, _parentctx, _parentState)
				p.PushNewRecursionContext(localctx, _startState, FqlParserRULE_expression)
				p.SetState(485)

				if !(p.Precpred(p.GetParserRuleContext(), 12)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 12)", ""))
				}
				{
					p.SetState(486)
					p.Match(FqlParserQuestionMark)
				}
				p.SetState(488)
				p.GetErrorHandler().Sync(p)
				_la = p.GetTokenStream().LA(1)

				if (((_la)&-(0x1f+1)) == 0 && ((1<<uint(_la))&((1<<FqlParserOpenBracket)|(1<<FqlParserOpenParen)|(1<<FqlParserOpenBrace)|(1<<FqlParserPlus)|(1<<FqlParserMinus))) != 0) || (((_la-44)&-(0x1f+1)) == 0 && ((1<<uint((_la-44)))&((1<<(FqlParserNone-44))|(1<<(FqlParserNull-44))|(1<<(FqlParserBooleanLiteral-44))|(1<<(FqlParserLike-44))|(1<<(FqlParserNot-44))|(1<<(FqlParserParam-44))|(1<<(FqlParserIdentifier-44))|(1<<(FqlParserStringLiteral-44))|(1<<(FqlParserTemplateStringLiteral-44))|(1<<(FqlParserIntegerLiteral-44))|(1<<(FqlParserFloatLiteral-44)))) != 0) {
					{
						p.SetState(487)
						p.expression(0)
					}

				}
				{
					p.SetState(490)
					p.Match(FqlParserColon)
				}
				{
					p.SetState(491)
					p.expression(13)
				}

			}

		}
		p.SetState(496)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 47, p.GetParserRuleContext())
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
	p.EnterRule(localctx, 92, FqlParserRULE_forTernaryExpression)
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

	p.SetState(525)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 49, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(497)
			p.expression(0)
		}
		{
			p.SetState(498)
			p.Match(FqlParserQuestionMark)
		}
		p.SetState(500)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		if (((_la)&-(0x1f+1)) == 0 && ((1<<uint(_la))&((1<<FqlParserOpenBracket)|(1<<FqlParserOpenParen)|(1<<FqlParserOpenBrace)|(1<<FqlParserPlus)|(1<<FqlParserMinus))) != 0) || (((_la-44)&-(0x1f+1)) == 0 && ((1<<uint((_la-44)))&((1<<(FqlParserNone-44))|(1<<(FqlParserNull-44))|(1<<(FqlParserBooleanLiteral-44))|(1<<(FqlParserLike-44))|(1<<(FqlParserNot-44))|(1<<(FqlParserParam-44))|(1<<(FqlParserIdentifier-44))|(1<<(FqlParserStringLiteral-44))|(1<<(FqlParserTemplateStringLiteral-44))|(1<<(FqlParserIntegerLiteral-44))|(1<<(FqlParserFloatLiteral-44)))) != 0) {
			{
				p.SetState(499)
				p.expression(0)
			}

		}
		{
			p.SetState(502)
			p.Match(FqlParserColon)
		}
		{
			p.SetState(503)
			p.Match(FqlParserOpenParen)
		}
		{
			p.SetState(504)
			p.ForExpression()
		}
		{
			p.SetState(505)
			p.Match(FqlParserCloseParen)
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(507)
			p.expression(0)
		}
		{
			p.SetState(508)
			p.Match(FqlParserQuestionMark)
		}
		{
			p.SetState(509)
			p.Match(FqlParserOpenParen)
		}
		{
			p.SetState(510)
			p.ForExpression()
		}
		{
			p.SetState(511)
			p.Match(FqlParserCloseParen)
		}
		{
			p.SetState(512)
			p.Match(FqlParserColon)
		}
		{
			p.SetState(513)
			p.expression(0)
		}

	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(515)
			p.expression(0)
		}
		{
			p.SetState(516)
			p.Match(FqlParserQuestionMark)
		}
		{
			p.SetState(517)
			p.Match(FqlParserOpenParen)
		}
		{
			p.SetState(518)
			p.ForExpression()
		}
		{
			p.SetState(519)
			p.Match(FqlParserCloseParen)
		}
		{
			p.SetState(520)
			p.Match(FqlParserColon)
		}
		{
			p.SetState(521)
			p.Match(FqlParserOpenParen)
		}
		{
			p.SetState(522)
			p.ForExpression()
		}
		{
			p.SetState(523)
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
	p.EnterRule(localctx, 94, FqlParserRULE_arrayOperator)
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
		p.SetState(527)
		_la = p.GetTokenStream().LA(1)

		if !(((_la-44)&-(0x1f+1)) == 0 && ((1<<uint((_la-44)))&((1<<(FqlParserNone-44))|(1<<(FqlParserAll-44))|(1<<(FqlParserAny-44)))) != 0) {
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
	p.EnterRule(localctx, 96, FqlParserRULE_inOperator)

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

	p.SetState(532)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case FqlParserIn:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(529)
			p.Match(FqlParserIn)
		}

	case FqlParserNot:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(530)
			p.Match(FqlParserNot)
		}
		{
			p.SetState(531)
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
	p.EnterRule(localctx, 98, FqlParserRULE_equalityOperator)
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
		p.SetState(534)
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

// ILogicalOperatorContext is an interface to support dynamic dispatch.
type ILogicalOperatorContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsLogicalOperatorContext differentiates from other interfaces.
	IsLogicalOperatorContext()
}

type LogicalOperatorContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyLogicalOperatorContext() *LogicalOperatorContext {
	var p = new(LogicalOperatorContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = FqlParserRULE_logicalOperator
	return p
}

func (*LogicalOperatorContext) IsLogicalOperatorContext() {}

func NewLogicalOperatorContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *LogicalOperatorContext {
	var p = new(LogicalOperatorContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_logicalOperator

	return p
}

func (s *LogicalOperatorContext) GetParser() antlr.Parser { return s.parser }

func (s *LogicalOperatorContext) And() antlr.TerminalNode {
	return s.GetToken(FqlParserAnd, 0)
}

func (s *LogicalOperatorContext) Or() antlr.TerminalNode {
	return s.GetToken(FqlParserOr, 0)
}

func (s *LogicalOperatorContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *LogicalOperatorContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *LogicalOperatorContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterLogicalOperator(s)
	}
}

func (s *LogicalOperatorContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitLogicalOperator(s)
	}
}

func (s *LogicalOperatorContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitLogicalOperator(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) LogicalOperator() (localctx ILogicalOperatorContext) {
	localctx = NewLogicalOperatorContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 100, FqlParserRULE_logicalOperator)
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
		p.SetState(536)
		_la = p.GetTokenStream().LA(1)

		if !(_la == FqlParserAnd || _la == FqlParserOr) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}

	return localctx
}

// IMathOperatorContext is an interface to support dynamic dispatch.
type IMathOperatorContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsMathOperatorContext differentiates from other interfaces.
	IsMathOperatorContext()
}

type MathOperatorContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyMathOperatorContext() *MathOperatorContext {
	var p = new(MathOperatorContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = FqlParserRULE_mathOperator
	return p
}

func (*MathOperatorContext) IsMathOperatorContext() {}

func NewMathOperatorContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *MathOperatorContext {
	var p = new(MathOperatorContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_mathOperator

	return p
}

func (s *MathOperatorContext) GetParser() antlr.Parser { return s.parser }

func (s *MathOperatorContext) Plus() antlr.TerminalNode {
	return s.GetToken(FqlParserPlus, 0)
}

func (s *MathOperatorContext) Minus() antlr.TerminalNode {
	return s.GetToken(FqlParserMinus, 0)
}

func (s *MathOperatorContext) Multi() antlr.TerminalNode {
	return s.GetToken(FqlParserMulti, 0)
}

func (s *MathOperatorContext) Div() antlr.TerminalNode {
	return s.GetToken(FqlParserDiv, 0)
}

func (s *MathOperatorContext) Mod() antlr.TerminalNode {
	return s.GetToken(FqlParserMod, 0)
}

func (s *MathOperatorContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *MathOperatorContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *MathOperatorContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterMathOperator(s)
	}
}

func (s *MathOperatorContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitMathOperator(s)
	}
}

func (s *MathOperatorContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitMathOperator(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) MathOperator() (localctx IMathOperatorContext) {
	localctx = NewMathOperatorContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 102, FqlParserRULE_mathOperator)
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
		p.SetState(538)
		_la = p.GetTokenStream().LA(1)

		if !(((_la)&-(0x1f+1)) == 0 && ((1<<uint(_la))&((1<<FqlParserPlus)|(1<<FqlParserMinus)|(1<<FqlParserMulti)|(1<<FqlParserDiv)|(1<<FqlParserMod))) != 0) {
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
	p.EnterRule(localctx, 104, FqlParserRULE_unaryOperator)
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
		p.SetState(540)
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
	case 45:
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
		return p.Precpred(p.GetParserRuleContext(), 19)

	case 1:
		return p.Precpred(p.GetParserRuleContext(), 18)

	case 2:
		return p.Precpred(p.GetParserRuleContext(), 17)

	case 3:
		return p.Precpred(p.GetParserRuleContext(), 14)

	case 4:
		return p.Precpred(p.GetParserRuleContext(), 13)

	case 5:
		return p.Precpred(p.GetParserRuleContext(), 12)

	default:
		panic("No predicate with index: " + fmt.Sprint(predIndex))
	}
}
