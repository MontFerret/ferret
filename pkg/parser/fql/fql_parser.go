// Code generated from antlr/FqlParser.g4 by ANTLR 4.8. DO NOT EDIT.

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
	3, 24715, 42794, 33075, 47597, 16764, 15335, 30598, 22884, 3, 65, 610,
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
	3, 2, 7, 2, 132, 10, 2, 12, 2, 14, 2, 135, 11, 2, 3, 2, 3, 2, 3, 3, 3,
	3, 3, 4, 3, 4, 3, 5, 3, 5, 3, 5, 3, 6, 7, 6, 147, 10, 6, 12, 6, 14, 6,
	150, 11, 6, 3, 6, 3, 6, 3, 7, 3, 7, 5, 7, 156, 10, 7, 3, 8, 3, 8, 5, 8,
	160, 10, 8, 3, 9, 3, 9, 5, 9, 164, 10, 9, 3, 9, 3, 9, 3, 9, 5, 9, 169,
	10, 9, 3, 9, 3, 9, 3, 9, 3, 9, 3, 9, 3, 9, 5, 9, 177, 10, 9, 3, 10, 3,
	10, 3, 10, 3, 10, 5, 10, 183, 10, 10, 3, 10, 3, 10, 3, 10, 7, 10, 188,
	10, 10, 12, 10, 14, 10, 191, 11, 10, 3, 10, 3, 10, 3, 11, 3, 11, 3, 12,
	3, 12, 3, 13, 3, 13, 3, 13, 3, 13, 3, 13, 3, 13, 3, 13, 5, 13, 206, 10,
	13, 3, 14, 3, 14, 3, 14, 3, 14, 5, 14, 212, 10, 14, 3, 15, 3, 15, 5, 15,
	216, 10, 15, 3, 16, 3, 16, 5, 16, 220, 10, 16, 3, 17, 3, 17, 5, 17, 224,
	10, 17, 3, 18, 3, 18, 3, 18, 3, 19, 3, 19, 3, 19, 3, 19, 5, 19, 233, 10,
	19, 3, 20, 3, 20, 5, 20, 237, 10, 20, 3, 21, 3, 21, 3, 21, 3, 21, 7, 21,
	243, 10, 21, 12, 21, 14, 21, 246, 11, 21, 3, 22, 3, 22, 5, 22, 250, 10,
	22, 3, 23, 3, 23, 3, 23, 3, 23, 3, 23, 3, 23, 3, 23, 3, 23, 3, 23, 3, 23,
	3, 23, 3, 23, 3, 23, 3, 23, 3, 23, 3, 23, 3, 23, 3, 23, 5, 23, 270, 10,
	23, 3, 24, 3, 24, 3, 24, 3, 24, 3, 25, 3, 25, 3, 25, 7, 25, 279, 10, 25,
	12, 25, 14, 25, 282, 11, 25, 3, 26, 3, 26, 3, 26, 3, 26, 7, 26, 288, 10,
	26, 12, 26, 14, 26, 291, 11, 26, 3, 27, 3, 27, 3, 27, 3, 27, 3, 28, 3,
	28, 3, 28, 3, 28, 3, 28, 3, 28, 5, 28, 303, 10, 28, 5, 28, 305, 10, 28,
	3, 29, 3, 29, 3, 29, 3, 29, 3, 29, 3, 30, 3, 30, 3, 30, 3, 30, 3, 30, 3,
	30, 3, 30, 3, 30, 3, 30, 3, 30, 3, 30, 3, 30, 3, 30, 3, 30, 3, 30, 5, 30,
	327, 10, 30, 3, 31, 3, 31, 3, 31, 3, 32, 3, 32, 3, 33, 3, 33, 3, 33, 5,
	33, 337, 10, 33, 3, 33, 3, 33, 3, 33, 3, 33, 5, 33, 343, 10, 33, 3, 34,
	3, 34, 5, 34, 347, 10, 34, 3, 34, 3, 34, 3, 35, 3, 35, 3, 35, 3, 35, 7,
	35, 355, 10, 35, 12, 35, 14, 35, 358, 11, 35, 5, 35, 360, 10, 35, 3, 35,
	5, 35, 363, 10, 35, 3, 35, 3, 35, 3, 36, 3, 36, 3, 37, 3, 37, 3, 38, 3,
	38, 3, 39, 3, 39, 3, 40, 3, 40, 3, 41, 3, 41, 6, 41, 379, 10, 41, 13, 41,
	14, 41, 380, 3, 41, 7, 41, 384, 10, 41, 12, 41, 14, 41, 387, 11, 41, 3,
	42, 3, 42, 3, 42, 3, 42, 3, 42, 3, 42, 3, 42, 3, 42, 3, 42, 5, 42, 398,
	10, 42, 3, 43, 3, 43, 3, 44, 3, 44, 3, 44, 3, 44, 3, 45, 3, 45, 3, 45,
	5, 45, 409, 10, 45, 3, 46, 3, 46, 3, 46, 3, 46, 3, 47, 3, 47, 3, 47, 3,
	48, 7, 48, 419, 10, 48, 12, 48, 14, 48, 422, 11, 48, 3, 49, 3, 49, 3, 50,
	3, 50, 3, 50, 3, 50, 3, 51, 3, 51, 3, 51, 5, 51, 433, 10, 51, 3, 52, 3,
	52, 3, 52, 7, 52, 438, 10, 52, 12, 52, 14, 52, 441, 11, 52, 6, 52, 443,
	10, 52, 13, 52, 14, 52, 444, 3, 52, 3, 52, 3, 52, 3, 52, 7, 52, 451, 10,
	52, 12, 52, 14, 52, 454, 11, 52, 7, 52, 456, 10, 52, 12, 52, 14, 52, 459,
	11, 52, 3, 52, 3, 52, 3, 52, 7, 52, 464, 10, 52, 12, 52, 14, 52, 467, 11,
	52, 7, 52, 469, 10, 52, 12, 52, 14, 52, 472, 11, 52, 5, 52, 474, 10, 52,
	3, 53, 3, 53, 3, 53, 3, 54, 3, 54, 3, 54, 3, 54, 7, 54, 483, 10, 54, 12,
	54, 14, 54, 486, 11, 54, 5, 54, 488, 10, 54, 3, 54, 3, 54, 3, 55, 3, 55,
	3, 55, 3, 55, 3, 55, 3, 55, 3, 55, 3, 55, 3, 55, 3, 55, 3, 55, 3, 55, 3,
	55, 3, 55, 3, 55, 3, 55, 3, 55, 5, 55, 509, 10, 55, 3, 55, 3, 55, 3, 55,
	3, 55, 3, 55, 3, 55, 3, 55, 3, 55, 3, 55, 3, 55, 3, 55, 3, 55, 5, 55, 523,
	10, 55, 3, 55, 3, 55, 3, 55, 3, 55, 3, 55, 3, 55, 3, 55, 3, 55, 3, 55,
	3, 55, 3, 55, 3, 55, 3, 55, 3, 55, 3, 55, 3, 55, 3, 55, 3, 55, 3, 55, 3,
	55, 3, 55, 3, 55, 3, 55, 3, 55, 3, 55, 5, 55, 550, 10, 55, 3, 55, 3, 55,
	7, 55, 554, 10, 55, 12, 55, 14, 55, 557, 11, 55, 3, 56, 3, 56, 3, 56, 5,
	56, 562, 10, 56, 3, 56, 3, 56, 3, 56, 3, 56, 3, 56, 3, 56, 3, 56, 3, 56,
	3, 56, 3, 56, 3, 56, 3, 56, 3, 56, 3, 56, 3, 56, 3, 56, 3, 56, 3, 56, 3,
	56, 3, 56, 3, 56, 3, 56, 3, 56, 5, 56, 587, 10, 56, 3, 57, 3, 57, 3, 58,
	3, 58, 3, 58, 5, 58, 594, 10, 58, 3, 59, 3, 59, 3, 60, 3, 60, 3, 61, 3,
	61, 3, 62, 3, 62, 3, 63, 3, 63, 3, 64, 3, 64, 3, 65, 3, 65, 3, 65, 2, 3,
	108, 66, 2, 4, 6, 8, 10, 12, 14, 16, 18, 20, 22, 24, 26, 28, 30, 32, 34,
	36, 38, 40, 42, 44, 46, 48, 50, 52, 54, 56, 58, 60, 62, 64, 66, 68, 70,
	72, 74, 76, 78, 80, 82, 84, 86, 88, 90, 92, 94, 96, 98, 100, 102, 104,
	106, 108, 110, 112, 114, 116, 118, 120, 122, 124, 126, 128, 2, 10, 3, 2,
	46, 47, 5, 2, 30, 31, 37, 59, 61, 61, 4, 2, 46, 46, 54, 55, 3, 2, 17, 22,
	3, 2, 35, 36, 3, 2, 23, 25, 3, 2, 26, 27, 4, 2, 26, 27, 57, 58, 2, 636,
	2, 133, 3, 2, 2, 2, 4, 138, 3, 2, 2, 2, 6, 140, 3, 2, 2, 2, 8, 142, 3,
	2, 2, 2, 10, 148, 3, 2, 2, 2, 12, 155, 3, 2, 2, 2, 14, 159, 3, 2, 2, 2,
	16, 176, 3, 2, 2, 2, 18, 178, 3, 2, 2, 2, 20, 194, 3, 2, 2, 2, 22, 196,
	3, 2, 2, 2, 24, 205, 3, 2, 2, 2, 26, 211, 3, 2, 2, 2, 28, 215, 3, 2, 2,
	2, 30, 219, 3, 2, 2, 2, 32, 223, 3, 2, 2, 2, 34, 225, 3, 2, 2, 2, 36, 228,
	3, 2, 2, 2, 38, 236, 3, 2, 2, 2, 40, 238, 3, 2, 2, 2, 42, 247, 3, 2, 2,
	2, 44, 269, 3, 2, 2, 2, 46, 271, 3, 2, 2, 2, 48, 275, 3, 2, 2, 2, 50, 283,
	3, 2, 2, 2, 52, 292, 3, 2, 2, 2, 54, 304, 3, 2, 2, 2, 56, 306, 3, 2, 2,
	2, 58, 326, 3, 2, 2, 2, 60, 328, 3, 2, 2, 2, 62, 331, 3, 2, 2, 2, 64, 336,
	3, 2, 2, 2, 66, 344, 3, 2, 2, 2, 68, 350, 3, 2, 2, 2, 70, 366, 3, 2, 2,
	2, 72, 368, 3, 2, 2, 2, 74, 370, 3, 2, 2, 2, 76, 372, 3, 2, 2, 2, 78, 374,
	3, 2, 2, 2, 80, 376, 3, 2, 2, 2, 82, 397, 3, 2, 2, 2, 84, 399, 3, 2, 2,
	2, 86, 401, 3, 2, 2, 2, 88, 408, 3, 2, 2, 2, 90, 410, 3, 2, 2, 2, 92, 414,
	3, 2, 2, 2, 94, 420, 3, 2, 2, 2, 96, 423, 3, 2, 2, 2, 98, 425, 3, 2, 2,
	2, 100, 432, 3, 2, 2, 2, 102, 473, 3, 2, 2, 2, 104, 475, 3, 2, 2, 2, 106,
	478, 3, 2, 2, 2, 108, 508, 3, 2, 2, 2, 110, 586, 3, 2, 2, 2, 112, 588,
	3, 2, 2, 2, 114, 593, 3, 2, 2, 2, 116, 595, 3, 2, 2, 2, 118, 597, 3, 2,
	2, 2, 120, 599, 3, 2, 2, 2, 122, 601, 3, 2, 2, 2, 124, 603, 3, 2, 2, 2,
	126, 605, 3, 2, 2, 2, 128, 607, 3, 2, 2, 2, 130, 132, 5, 4, 3, 2, 131,
	130, 3, 2, 2, 2, 132, 135, 3, 2, 2, 2, 133, 131, 3, 2, 2, 2, 133, 134,
	3, 2, 2, 2, 134, 136, 3, 2, 2, 2, 135, 133, 3, 2, 2, 2, 136, 137, 5, 10,
	6, 2, 137, 3, 3, 2, 2, 2, 138, 139, 5, 6, 4, 2, 139, 5, 3, 2, 2, 2, 140,
	141, 5, 8, 5, 2, 141, 7, 3, 2, 2, 2, 142, 143, 7, 49, 2, 2, 143, 144, 5,
	92, 47, 2, 144, 9, 3, 2, 2, 2, 145, 147, 5, 12, 7, 2, 146, 145, 3, 2, 2,
	2, 147, 150, 3, 2, 2, 2, 148, 146, 3, 2, 2, 2, 148, 149, 3, 2, 2, 2, 149,
	151, 3, 2, 2, 2, 150, 148, 3, 2, 2, 2, 151, 152, 5, 14, 8, 2, 152, 11,
	3, 2, 2, 2, 153, 156, 5, 98, 50, 2, 154, 156, 5, 58, 30, 2, 155, 153, 3,
	2, 2, 2, 155, 154, 3, 2, 2, 2, 156, 13, 3, 2, 2, 2, 157, 160, 5, 16, 9,
	2, 158, 160, 5, 18, 10, 2, 159, 157, 3, 2, 2, 2, 159, 158, 3, 2, 2, 2,
	160, 15, 3, 2, 2, 2, 161, 163, 7, 38, 2, 2, 162, 164, 7, 39, 2, 2, 163,
	162, 3, 2, 2, 2, 163, 164, 3, 2, 2, 2, 164, 165, 3, 2, 2, 2, 165, 177,
	5, 108, 55, 2, 166, 168, 7, 38, 2, 2, 167, 169, 7, 39, 2, 2, 168, 167,
	3, 2, 2, 2, 168, 169, 3, 2, 2, 2, 169, 170, 3, 2, 2, 2, 170, 171, 7, 13,
	2, 2, 171, 172, 5, 18, 10, 2, 172, 173, 7, 14, 2, 2, 173, 177, 3, 2, 2,
	2, 174, 175, 7, 38, 2, 2, 175, 177, 5, 110, 56, 2, 176, 161, 3, 2, 2, 2,
	176, 166, 3, 2, 2, 2, 176, 174, 3, 2, 2, 2, 177, 17, 3, 2, 2, 2, 178, 179,
	7, 37, 2, 2, 179, 182, 5, 20, 11, 2, 180, 181, 7, 10, 2, 2, 181, 183, 5,
	22, 12, 2, 182, 180, 3, 2, 2, 2, 182, 183, 3, 2, 2, 2, 183, 184, 3, 2,
	2, 2, 184, 185, 7, 59, 2, 2, 185, 189, 5, 24, 13, 2, 186, 188, 5, 30, 16,
	2, 187, 186, 3, 2, 2, 2, 188, 191, 3, 2, 2, 2, 189, 187, 3, 2, 2, 2, 189,
	190, 3, 2, 2, 2, 190, 192, 3, 2, 2, 2, 191, 189, 3, 2, 2, 2, 192, 193,
	5, 32, 17, 2, 193, 19, 3, 2, 2, 2, 194, 195, 7, 61, 2, 2, 195, 21, 3, 2,
	2, 2, 196, 197, 7, 61, 2, 2, 197, 23, 3, 2, 2, 2, 198, 206, 5, 98, 50,
	2, 199, 206, 5, 66, 34, 2, 200, 206, 5, 68, 35, 2, 201, 206, 5, 62, 32,
	2, 202, 206, 5, 104, 53, 2, 203, 206, 5, 64, 33, 2, 204, 206, 5, 60, 31,
	2, 205, 198, 3, 2, 2, 2, 205, 199, 3, 2, 2, 2, 205, 200, 3, 2, 2, 2, 205,
	201, 3, 2, 2, 2, 205, 202, 3, 2, 2, 2, 205, 203, 3, 2, 2, 2, 205, 204,
	3, 2, 2, 2, 206, 25, 3, 2, 2, 2, 207, 212, 5, 36, 19, 2, 208, 212, 5, 40,
	21, 2, 209, 212, 5, 34, 18, 2, 210, 212, 5, 44, 23, 2, 211, 207, 3, 2,
	2, 2, 211, 208, 3, 2, 2, 2, 211, 209, 3, 2, 2, 2, 211, 210, 3, 2, 2, 2,
	212, 27, 3, 2, 2, 2, 213, 216, 5, 58, 30, 2, 214, 216, 5, 98, 50, 2, 215,
	213, 3, 2, 2, 2, 215, 214, 3, 2, 2, 2, 216, 29, 3, 2, 2, 2, 217, 220, 5,
	28, 15, 2, 218, 220, 5, 26, 14, 2, 219, 217, 3, 2, 2, 2, 219, 218, 3, 2,
	2, 2, 220, 31, 3, 2, 2, 2, 221, 224, 5, 16, 9, 2, 222, 224, 5, 18, 10,
	2, 223, 221, 3, 2, 2, 2, 223, 222, 3, 2, 2, 2, 224, 33, 3, 2, 2, 2, 225,
	226, 7, 40, 2, 2, 226, 227, 5, 108, 55, 2, 227, 35, 3, 2, 2, 2, 228, 229,
	7, 42, 2, 2, 229, 232, 5, 38, 20, 2, 230, 231, 7, 10, 2, 2, 231, 233, 5,
	38, 20, 2, 232, 230, 3, 2, 2, 2, 232, 233, 3, 2, 2, 2, 233, 37, 3, 2, 2,
	2, 234, 237, 7, 63, 2, 2, 235, 237, 5, 60, 31, 2, 236, 234, 3, 2, 2, 2,
	236, 235, 3, 2, 2, 2, 237, 39, 3, 2, 2, 2, 238, 239, 7, 41, 2, 2, 239,
	244, 5, 42, 22, 2, 240, 241, 7, 10, 2, 2, 241, 243, 5, 42, 22, 2, 242,
	240, 3, 2, 2, 2, 243, 246, 3, 2, 2, 2, 244, 242, 3, 2, 2, 2, 244, 245,
	3, 2, 2, 2, 245, 41, 3, 2, 2, 2, 246, 244, 3, 2, 2, 2, 247, 249, 5, 108,
	55, 2, 248, 250, 7, 45, 2, 2, 249, 248, 3, 2, 2, 2, 249, 250, 3, 2, 2,
	2, 250, 43, 3, 2, 2, 2, 251, 252, 7, 44, 2, 2, 252, 270, 5, 56, 29, 2,
	253, 254, 7, 44, 2, 2, 254, 270, 5, 50, 26, 2, 255, 256, 7, 44, 2, 2, 256,
	257, 5, 48, 25, 2, 257, 258, 5, 50, 26, 2, 258, 270, 3, 2, 2, 2, 259, 260,
	7, 44, 2, 2, 260, 261, 5, 48, 25, 2, 261, 262, 5, 54, 28, 2, 262, 270,
	3, 2, 2, 2, 263, 264, 7, 44, 2, 2, 264, 265, 5, 48, 25, 2, 265, 266, 5,
	56, 29, 2, 266, 270, 3, 2, 2, 2, 267, 268, 7, 44, 2, 2, 268, 270, 5, 48,
	25, 2, 269, 251, 3, 2, 2, 2, 269, 253, 3, 2, 2, 2, 269, 255, 3, 2, 2, 2,
	269, 259, 3, 2, 2, 2, 269, 263, 3, 2, 2, 2, 269, 267, 3, 2, 2, 2, 270,
	45, 3, 2, 2, 2, 271, 272, 7, 61, 2, 2, 272, 273, 7, 33, 2, 2, 273, 274,
	5, 108, 55, 2, 274, 47, 3, 2, 2, 2, 275, 280, 5, 46, 24, 2, 276, 277, 7,
	10, 2, 2, 277, 279, 5, 46, 24, 2, 278, 276, 3, 2, 2, 2, 279, 282, 3, 2,
	2, 2, 280, 278, 3, 2, 2, 2, 280, 281, 3, 2, 2, 2, 281, 49, 3, 2, 2, 2,
	282, 280, 3, 2, 2, 2, 283, 284, 7, 56, 2, 2, 284, 289, 5, 52, 27, 2, 285,
	286, 7, 10, 2, 2, 286, 288, 5, 52, 27, 2, 287, 285, 3, 2, 2, 2, 288, 291,
	3, 2, 2, 2, 289, 287, 3, 2, 2, 2, 289, 290, 3, 2, 2, 2, 290, 51, 3, 2,
	2, 2, 291, 289, 3, 2, 2, 2, 292, 293, 7, 61, 2, 2, 293, 294, 7, 33, 2,
	2, 294, 295, 5, 98, 50, 2, 295, 53, 3, 2, 2, 2, 296, 297, 7, 50, 2, 2,
	297, 305, 5, 46, 24, 2, 298, 299, 7, 50, 2, 2, 299, 302, 7, 61, 2, 2, 300,
	301, 7, 51, 2, 2, 301, 303, 7, 61, 2, 2, 302, 300, 3, 2, 2, 2, 302, 303,
	3, 2, 2, 2, 303, 305, 3, 2, 2, 2, 304, 296, 3, 2, 2, 2, 304, 298, 3, 2,
	2, 2, 305, 55, 3, 2, 2, 2, 306, 307, 7, 52, 2, 2, 307, 308, 7, 53, 2, 2,
	308, 309, 7, 50, 2, 2, 309, 310, 7, 61, 2, 2, 310, 57, 3, 2, 2, 2, 311,
	312, 7, 43, 2, 2, 312, 313, 7, 61, 2, 2, 313, 314, 7, 33, 2, 2, 314, 327,
	5, 108, 55, 2, 315, 316, 7, 43, 2, 2, 316, 317, 7, 61, 2, 2, 317, 318,
	7, 33, 2, 2, 318, 319, 7, 13, 2, 2, 319, 320, 5, 18, 10, 2, 320, 321, 7,
	14, 2, 2, 321, 327, 3, 2, 2, 2, 322, 323, 7, 43, 2, 2, 323, 324, 7, 61,
	2, 2, 324, 325, 7, 33, 2, 2, 325, 327, 5, 110, 56, 2, 326, 311, 3, 2, 2,
	2, 326, 315, 3, 2, 2, 2, 326, 322, 3, 2, 2, 2, 327, 59, 3, 2, 2, 2, 328,
	329, 7, 60, 2, 2, 329, 330, 7, 61, 2, 2, 330, 61, 3, 2, 2, 2, 331, 332,
	7, 61, 2, 2, 332, 63, 3, 2, 2, 2, 333, 337, 5, 74, 38, 2, 334, 337, 5,
	62, 32, 2, 335, 337, 5, 60, 31, 2, 336, 333, 3, 2, 2, 2, 336, 334, 3, 2,
	2, 2, 336, 335, 3, 2, 2, 2, 337, 338, 3, 2, 2, 2, 338, 342, 7, 32, 2, 2,
	339, 343, 5, 74, 38, 2, 340, 343, 5, 62, 32, 2, 341, 343, 5, 60, 31, 2,
	342, 339, 3, 2, 2, 2, 342, 340, 3, 2, 2, 2, 342, 341, 3, 2, 2, 2, 343,
	65, 3, 2, 2, 2, 344, 346, 7, 11, 2, 2, 345, 347, 5, 80, 41, 2, 346, 345,
	3, 2, 2, 2, 346, 347, 3, 2, 2, 2, 347, 348, 3, 2, 2, 2, 348, 349, 7, 12,
	2, 2, 349, 67, 3, 2, 2, 2, 350, 359, 7, 15, 2, 2, 351, 356, 5, 82, 42,
	2, 352, 353, 7, 10, 2, 2, 353, 355, 5, 82, 42, 2, 354, 352, 3, 2, 2, 2,
	355, 358, 3, 2, 2, 2, 356, 354, 3, 2, 2, 2, 356, 357, 3, 2, 2, 2, 357,
	360, 3, 2, 2, 2, 358, 356, 3, 2, 2, 2, 359, 351, 3, 2, 2, 2, 359, 360,
	3, 2, 2, 2, 360, 362, 3, 2, 2, 2, 361, 363, 7, 10, 2, 2, 362, 361, 3, 2,
	2, 2, 362, 363, 3, 2, 2, 2, 363, 364, 3, 2, 2, 2, 364, 365, 7, 16, 2, 2,
	365, 69, 3, 2, 2, 2, 366, 367, 7, 48, 2, 2, 367, 71, 3, 2, 2, 2, 368, 369,
	7, 62, 2, 2, 369, 73, 3, 2, 2, 2, 370, 371, 7, 63, 2, 2, 371, 75, 3, 2,
	2, 2, 372, 373, 7, 64, 2, 2, 373, 77, 3, 2, 2, 2, 374, 375, 9, 2, 2, 2,
	375, 79, 3, 2, 2, 2, 376, 385, 5, 108, 55, 2, 377, 379, 7, 10, 2, 2, 378,
	377, 3, 2, 2, 2, 379, 380, 3, 2, 2, 2, 380, 378, 3, 2, 2, 2, 380, 381,
	3, 2, 2, 2, 381, 382, 3, 2, 2, 2, 382, 384, 5, 108, 55, 2, 383, 378, 3,
	2, 2, 2, 384, 387, 3, 2, 2, 2, 385, 383, 3, 2, 2, 2, 385, 386, 3, 2, 2,
	2, 386, 81, 3, 2, 2, 2, 387, 385, 3, 2, 2, 2, 388, 389, 5, 88, 45, 2, 389,
	390, 7, 7, 2, 2, 390, 391, 5, 108, 55, 2, 391, 398, 3, 2, 2, 2, 392, 393,
	5, 86, 44, 2, 393, 394, 7, 7, 2, 2, 394, 395, 5, 108, 55, 2, 395, 398,
	3, 2, 2, 2, 396, 398, 5, 84, 43, 2, 397, 388, 3, 2, 2, 2, 397, 392, 3,
	2, 2, 2, 397, 396, 3, 2, 2, 2, 398, 83, 3, 2, 2, 2, 399, 400, 5, 62, 32,
	2, 400, 85, 3, 2, 2, 2, 401, 402, 7, 11, 2, 2, 402, 403, 5, 108, 55, 2,
	403, 404, 7, 12, 2, 2, 404, 87, 3, 2, 2, 2, 405, 409, 7, 61, 2, 2, 406,
	409, 5, 72, 37, 2, 407, 409, 5, 60, 31, 2, 408, 405, 3, 2, 2, 2, 408, 406,
	3, 2, 2, 2, 408, 407, 3, 2, 2, 2, 409, 89, 3, 2, 2, 2, 410, 411, 7, 13,
	2, 2, 411, 412, 5, 108, 55, 2, 412, 413, 7, 14, 2, 2, 413, 91, 3, 2, 2,
	2, 414, 415, 5, 94, 48, 2, 415, 416, 7, 61, 2, 2, 416, 93, 3, 2, 2, 2,
	417, 419, 7, 65, 2, 2, 418, 417, 3, 2, 2, 2, 419, 422, 3, 2, 2, 2, 420,
	418, 3, 2, 2, 2, 420, 421, 3, 2, 2, 2, 421, 95, 3, 2, 2, 2, 422, 420, 3,
	2, 2, 2, 423, 424, 9, 3, 2, 2, 424, 97, 3, 2, 2, 2, 425, 426, 5, 94, 48,
	2, 426, 427, 5, 96, 49, 2, 427, 428, 5, 106, 54, 2, 428, 99, 3, 2, 2, 2,
	429, 433, 7, 61, 2, 2, 430, 433, 5, 98, 50, 2, 431, 433, 5, 60, 31, 2,
	432, 429, 3, 2, 2, 2, 432, 430, 3, 2, 2, 2, 432, 431, 3, 2, 2, 2, 433,
	101, 3, 2, 2, 2, 434, 435, 7, 9, 2, 2, 435, 439, 5, 88, 45, 2, 436, 438,
	5, 86, 44, 2, 437, 436, 3, 2, 2, 2, 438, 441, 3, 2, 2, 2, 439, 437, 3,
	2, 2, 2, 439, 440, 3, 2, 2, 2, 440, 443, 3, 2, 2, 2, 441, 439, 3, 2, 2,
	2, 442, 434, 3, 2, 2, 2, 443, 444, 3, 2, 2, 2, 444, 442, 3, 2, 2, 2, 444,
	445, 3, 2, 2, 2, 445, 474, 3, 2, 2, 2, 446, 457, 5, 86, 44, 2, 447, 448,
	7, 9, 2, 2, 448, 452, 5, 88, 45, 2, 449, 451, 5, 86, 44, 2, 450, 449, 3,
	2, 2, 2, 451, 454, 3, 2, 2, 2, 452, 450, 3, 2, 2, 2, 452, 453, 3, 2, 2,
	2, 453, 456, 3, 2, 2, 2, 454, 452, 3, 2, 2, 2, 455, 447, 3, 2, 2, 2, 456,
	459, 3, 2, 2, 2, 457, 455, 3, 2, 2, 2, 457, 458, 3, 2, 2, 2, 458, 470,
	3, 2, 2, 2, 459, 457, 3, 2, 2, 2, 460, 465, 5, 86, 44, 2, 461, 462, 7,
	9, 2, 2, 462, 464, 5, 88, 45, 2, 463, 461, 3, 2, 2, 2, 464, 467, 3, 2,
	2, 2, 465, 463, 3, 2, 2, 2, 465, 466, 3, 2, 2, 2, 466, 469, 3, 2, 2, 2,
	467, 465, 3, 2, 2, 2, 468, 460, 3, 2, 2, 2, 469, 472, 3, 2, 2, 2, 470,
	468, 3, 2, 2, 2, 470, 471, 3, 2, 2, 2, 471, 474, 3, 2, 2, 2, 472, 470,
	3, 2, 2, 2, 473, 442, 3, 2, 2, 2, 473, 446, 3, 2, 2, 2, 474, 103, 3, 2,
	2, 2, 475, 476, 5, 100, 51, 2, 476, 477, 5, 102, 52, 2, 477, 105, 3, 2,
	2, 2, 478, 487, 7, 13, 2, 2, 479, 484, 5, 108, 55, 2, 480, 481, 7, 10,
	2, 2, 481, 483, 5, 108, 55, 2, 482, 480, 3, 2, 2, 2, 483, 486, 3, 2, 2,
	2, 484, 482, 3, 2, 2, 2, 484, 485, 3, 2, 2, 2, 485, 488, 3, 2, 2, 2, 486,
	484, 3, 2, 2, 2, 487, 479, 3, 2, 2, 2, 487, 488, 3, 2, 2, 2, 488, 489,
	3, 2, 2, 2, 489, 490, 7, 14, 2, 2, 490, 107, 3, 2, 2, 2, 491, 492, 8, 55,
	1, 2, 492, 493, 5, 128, 65, 2, 493, 494, 5, 108, 55, 25, 494, 509, 3, 2,
	2, 2, 495, 509, 5, 98, 50, 2, 496, 509, 5, 90, 46, 2, 497, 509, 5, 64,
	33, 2, 498, 509, 5, 72, 37, 2, 499, 509, 5, 74, 38, 2, 500, 509, 5, 76,
	39, 2, 501, 509, 5, 70, 36, 2, 502, 509, 5, 66, 34, 2, 503, 509, 5, 68,
	35, 2, 504, 509, 5, 62, 32, 2, 505, 509, 5, 104, 53, 2, 506, 509, 5, 78,
	40, 2, 507, 509, 5, 60, 31, 2, 508, 491, 3, 2, 2, 2, 508, 495, 3, 2, 2,
	2, 508, 496, 3, 2, 2, 2, 508, 497, 3, 2, 2, 2, 508, 498, 3, 2, 2, 2, 508,
	499, 3, 2, 2, 2, 508, 500, 3, 2, 2, 2, 508, 501, 3, 2, 2, 2, 508, 502,
	3, 2, 2, 2, 508, 503, 3, 2, 2, 2, 508, 504, 3, 2, 2, 2, 508, 505, 3, 2,
	2, 2, 508, 506, 3, 2, 2, 2, 508, 507, 3, 2, 2, 2, 509, 555, 3, 2, 2, 2,
	510, 511, 12, 24, 2, 2, 511, 512, 5, 124, 63, 2, 512, 513, 5, 108, 55,
	25, 513, 554, 3, 2, 2, 2, 514, 515, 12, 23, 2, 2, 515, 516, 5, 126, 64,
	2, 516, 517, 5, 108, 55, 24, 517, 554, 3, 2, 2, 2, 518, 519, 12, 20, 2,
	2, 519, 522, 5, 112, 57, 2, 520, 523, 5, 114, 58, 2, 521, 523, 5, 116,
	59, 2, 522, 520, 3, 2, 2, 2, 522, 521, 3, 2, 2, 2, 523, 524, 3, 2, 2, 2,
	524, 525, 5, 108, 55, 21, 525, 554, 3, 2, 2, 2, 526, 527, 12, 19, 2, 2,
	527, 528, 5, 114, 58, 2, 528, 529, 5, 108, 55, 20, 529, 554, 3, 2, 2, 2,
	530, 531, 12, 18, 2, 2, 531, 532, 5, 116, 59, 2, 532, 533, 5, 108, 55,
	19, 533, 554, 3, 2, 2, 2, 534, 535, 12, 17, 2, 2, 535, 536, 5, 118, 60,
	2, 536, 537, 5, 108, 55, 18, 537, 554, 3, 2, 2, 2, 538, 539, 12, 16, 2,
	2, 539, 540, 5, 120, 61, 2, 540, 541, 5, 108, 55, 17, 541, 554, 3, 2, 2,
	2, 542, 543, 12, 15, 2, 2, 543, 544, 5, 122, 62, 2, 544, 545, 5, 108, 55,
	16, 545, 554, 3, 2, 2, 2, 546, 547, 12, 14, 2, 2, 547, 549, 7, 34, 2, 2,
	548, 550, 5, 108, 55, 2, 549, 548, 3, 2, 2, 2, 549, 550, 3, 2, 2, 2, 550,
	551, 3, 2, 2, 2, 551, 552, 7, 7, 2, 2, 552, 554, 5, 108, 55, 15, 553, 510,
	3, 2, 2, 2, 553, 514, 3, 2, 2, 2, 553, 518, 3, 2, 2, 2, 553, 526, 3, 2,
	2, 2, 553, 530, 3, 2, 2, 2, 553, 534, 3, 2, 2, 2, 553, 538, 3, 2, 2, 2,
	553, 542, 3, 2, 2, 2, 553, 546, 3, 2, 2, 2, 554, 557, 3, 2, 2, 2, 555,
	553, 3, 2, 2, 2, 555, 556, 3, 2, 2, 2, 556, 109, 3, 2, 2, 2, 557, 555,
	3, 2, 2, 2, 558, 559, 5, 108, 55, 2, 559, 561, 7, 34, 2, 2, 560, 562, 5,
	108, 55, 2, 561, 560, 3, 2, 2, 2, 561, 562, 3, 2, 2, 2, 562, 563, 3, 2,
	2, 2, 563, 564, 7, 7, 2, 2, 564, 565, 7, 13, 2, 2, 565, 566, 5, 18, 10,
	2, 566, 567, 7, 14, 2, 2, 567, 587, 3, 2, 2, 2, 568, 569, 5, 108, 55, 2,
	569, 570, 7, 34, 2, 2, 570, 571, 7, 13, 2, 2, 571, 572, 5, 18, 10, 2, 572,
	573, 7, 14, 2, 2, 573, 574, 7, 7, 2, 2, 574, 575, 5, 108, 55, 2, 575, 587,
	3, 2, 2, 2, 576, 577, 5, 108, 55, 2, 577, 578, 7, 34, 2, 2, 578, 579, 7,
	13, 2, 2, 579, 580, 5, 18, 10, 2, 580, 581, 7, 14, 2, 2, 581, 582, 7, 7,
	2, 2, 582, 583, 7, 13, 2, 2, 583, 584, 5, 18, 10, 2, 584, 585, 7, 14, 2,
	2, 585, 587, 3, 2, 2, 2, 586, 558, 3, 2, 2, 2, 586, 568, 3, 2, 2, 2, 586,
	576, 3, 2, 2, 2, 587, 111, 3, 2, 2, 2, 588, 589, 9, 4, 2, 2, 589, 113,
	3, 2, 2, 2, 590, 594, 7, 59, 2, 2, 591, 592, 7, 58, 2, 2, 592, 594, 7,
	59, 2, 2, 593, 590, 3, 2, 2, 2, 593, 591, 3, 2, 2, 2, 594, 115, 3, 2, 2,
	2, 595, 596, 9, 5, 2, 2, 596, 117, 3, 2, 2, 2, 597, 598, 9, 6, 2, 2, 598,
	119, 3, 2, 2, 2, 599, 600, 7, 30, 2, 2, 600, 121, 3, 2, 2, 2, 601, 602,
	7, 31, 2, 2, 602, 123, 3, 2, 2, 2, 603, 604, 9, 7, 2, 2, 604, 125, 3, 2,
	2, 2, 605, 606, 9, 8, 2, 2, 606, 127, 3, 2, 2, 2, 607, 608, 9, 9, 2, 2,
	608, 129, 3, 2, 2, 2, 55, 133, 148, 155, 159, 163, 168, 176, 182, 189,
	205, 211, 215, 219, 223, 232, 236, 244, 249, 269, 280, 289, 302, 304, 326,
	336, 342, 346, 356, 359, 362, 380, 385, 397, 408, 420, 432, 439, 444, 452,
	457, 465, 470, 473, 484, 487, 508, 522, 549, 553, 555, 561, 586, 593,
}
var deserializer = antlr.NewATNDeserializer(nil)
var deserializedATN = deserializer.DeserializeFromUInt16(parserATN)

var literalNames = []string{
	"", "", "", "", "", "':'", "';'", "'.'", "','", "'['", "']'", "'('", "')'",
	"'{'", "'}'", "'>'", "'<'", "'=='", "'>='", "'<='", "'!='", "'*'", "'/'",
	"'%'", "'+'", "'-'", "'--'", "'++'", "", "", "", "'='", "'?'", "'!~'",
	"'=~'", "'FOR'", "'RETURN'", "'DISTINCT'", "'FILTER'", "'SORT'", "'LIMIT'",
	"'LET'", "'COLLECT'", "", "'NONE'", "'NULL'", "", "'USE'", "'INTO'", "'KEEP'",
	"'WITH'", "'COUNT'", "'ALL'", "'ANY'", "'AGGREGATE'", "'LIKE'", "", "'IN'",
	"'@'",
}
var symbolicNames = []string{
	"", "MultiLineComment", "SingleLineComment", "WhiteSpaces", "LineTerminator",
	"Colon", "SemiColon", "Dot", "Comma", "OpenBracket", "CloseBracket", "OpenParen",
	"CloseParen", "OpenBrace", "CloseBrace", "Gt", "Lt", "Eq", "Gte", "Lte",
	"Neq", "Multi", "Div", "Mod", "Plus", "Minus", "MinusMinus", "PlusPlus",
	"And", "Or", "Range", "Assign", "QuestionMark", "RegexNotMatch", "RegexMatch",
	"For", "Return", "Distinct", "Filter", "Sort", "Limit", "Let", "Collect",
	"SortDirection", "None", "Null", "BooleanLiteral", "Use", "Into", "Keep",
	"With", "Count", "All", "Any", "Aggregate", "Like", "Not", "In", "Param",
	"Identifier", "StringLiteral", "IntegerLiteral", "FloatLiteral", "NamespaceSegment",
}

var ruleNames = []string{
	"program", "head", "useExpression", "use", "body", "bodyStatement", "bodyExpression",
	"returnExpression", "forExpression", "forExpressionValueVariable", "forExpressionKeyVariable",
	"forExpressionSource", "forExpressionClause", "forExpressionStatement",
	"forExpressionBody", "forExpressionReturn", "filterClause", "limitClause",
	"limitClauseValue", "sortClause", "sortClauseExpression", "collectClause",
	"collectSelector", "collectGrouping", "collectAggregator", "collectAggregateSelector",
	"collectGroupVariable", "collectCounter", "variableDeclaration", "param",
	"variable", "rangeOperator", "arrayLiteral", "objectLiteral", "booleanLiteral",
	"stringLiteral", "integerLiteral", "floatLiteral", "noneLiteral", "arrayElementList",
	"propertyAssignment", "shorthandPropertyName", "computedPropertyName",
	"propertyName", "expressionGroup", "namespaceIdentifier", "namespace",
	"functionIdentifier", "functionCallExpression", "member", "memberPath",
	"memberExpression", "arguments", "expression", "forTernaryExpression",
	"arrayOperator", "inOperator", "equalityOperator", "regexpOperator", "logicalAndOperator",
	"logicalOrOperator", "multiplicativeOperator", "additiveOperator", "unaryOperator",
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
	FqlParserDistinct          = 37
	FqlParserFilter            = 38
	FqlParserSort              = 39
	FqlParserLimit             = 40
	FqlParserLet               = 41
	FqlParserCollect           = 42
	FqlParserSortDirection     = 43
	FqlParserNone              = 44
	FqlParserNull              = 45
	FqlParserBooleanLiteral    = 46
	FqlParserUse               = 47
	FqlParserInto              = 48
	FqlParserKeep              = 49
	FqlParserWith              = 50
	FqlParserCount             = 51
	FqlParserAll               = 52
	FqlParserAny               = 53
	FqlParserAggregate         = 54
	FqlParserLike              = 55
	FqlParserNot               = 56
	FqlParserIn                = 57
	FqlParserParam             = 58
	FqlParserIdentifier        = 59
	FqlParserStringLiteral     = 60
	FqlParserIntegerLiteral    = 61
	FqlParserFloatLiteral      = 62
	FqlParserNamespaceSegment  = 63
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
	FqlParserRULE_variableDeclaration        = 28
	FqlParserRULE_param                      = 29
	FqlParserRULE_variable                   = 30
	FqlParserRULE_rangeOperator              = 31
	FqlParserRULE_arrayLiteral               = 32
	FqlParserRULE_objectLiteral              = 33
	FqlParserRULE_booleanLiteral             = 34
	FqlParserRULE_stringLiteral              = 35
	FqlParserRULE_integerLiteral             = 36
	FqlParserRULE_floatLiteral               = 37
	FqlParserRULE_noneLiteral                = 38
	FqlParserRULE_arrayElementList           = 39
	FqlParserRULE_propertyAssignment         = 40
	FqlParserRULE_shorthandPropertyName      = 41
	FqlParserRULE_computedPropertyName       = 42
	FqlParserRULE_propertyName               = 43
	FqlParserRULE_expressionGroup            = 44
	FqlParserRULE_namespaceIdentifier        = 45
	FqlParserRULE_namespace                  = 46
	FqlParserRULE_functionIdentifier         = 47
	FqlParserRULE_functionCallExpression     = 48
	FqlParserRULE_member                     = 49
	FqlParserRULE_memberPath                 = 50
	FqlParserRULE_memberExpression           = 51
	FqlParserRULE_arguments                  = 52
	FqlParserRULE_expression                 = 53
	FqlParserRULE_forTernaryExpression       = 54
	FqlParserRULE_arrayOperator              = 55
	FqlParserRULE_inOperator                 = 56
	FqlParserRULE_equalityOperator           = 57
	FqlParserRULE_regexpOperator             = 58
	FqlParserRULE_logicalAndOperator         = 59
	FqlParserRULE_logicalOrOperator          = 60
	FqlParserRULE_multiplicativeOperator     = 61
	FqlParserRULE_additiveOperator           = 62
	FqlParserRULE_unaryOperator              = 63
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
	p.SetState(131)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 0, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			{
				p.SetState(128)
				p.Head()
			}

		}
		p.SetState(133)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 0, p.GetParserRuleContext())
	}
	{
		p.SetState(134)
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
		p.SetState(136)
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
		p.SetState(138)
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
		p.SetState(140)
		p.Match(FqlParserUse)
	}
	{
		p.SetState(141)
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
	p.SetState(146)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 1, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			{
				p.SetState(143)
				p.BodyStatement()
			}

		}
		p.SetState(148)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 1, p.GetParserRuleContext())
	}
	{
		p.SetState(149)
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

	p.SetState(153)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 2, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(151)
			p.FunctionCallExpression()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(152)
			p.VariableDeclaration()
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

	p.SetState(157)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case FqlParserReturn:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(155)
			p.ReturnExpression()
		}

	case FqlParserFor:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(156)
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

	p.SetState(174)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 6, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(159)
			p.Match(FqlParserReturn)
		}
		p.SetState(161)
		p.GetErrorHandler().Sync(p)

		if p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 4, p.GetParserRuleContext()) == 1 {
			{
				p.SetState(160)
				p.Match(FqlParserDistinct)
			}

		}
		{
			p.SetState(163)
			p.expression(0)
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(164)
			p.Match(FqlParserReturn)
		}
		p.SetState(166)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		if _la == FqlParserDistinct {
			{
				p.SetState(165)
				p.Match(FqlParserDistinct)
			}

		}
		{
			p.SetState(168)
			p.Match(FqlParserOpenParen)
		}
		{
			p.SetState(169)
			p.ForExpression()
		}
		{
			p.SetState(170)
			p.Match(FqlParserCloseParen)
		}

	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(172)
			p.Match(FqlParserReturn)
		}
		{
			p.SetState(173)
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

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(176)
		p.Match(FqlParserFor)
	}
	{
		p.SetState(177)
		p.ForExpressionValueVariable()
	}
	p.SetState(180)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == FqlParserComma {
		{
			p.SetState(178)
			p.Match(FqlParserComma)
		}
		{
			p.SetState(179)
			p.ForExpressionKeyVariable()
		}

	}
	{
		p.SetState(182)
		p.Match(FqlParserIn)
	}
	{
		p.SetState(183)
		p.ForExpressionSource()
	}
	p.SetState(187)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 8, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			{
				p.SetState(184)
				p.ForExpressionBody()
			}

		}
		p.SetState(189)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 8, p.GetParserRuleContext())
	}
	{
		p.SetState(190)
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
		p.SetState(192)
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
		p.SetState(194)
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

	p.SetState(203)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 9, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(196)
			p.FunctionCallExpression()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(197)
			p.ArrayLiteral()
		}

	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(198)
			p.ObjectLiteral()
		}

	case 4:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(199)
			p.Variable()
		}

	case 5:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(200)
			p.MemberExpression()
		}

	case 6:
		p.EnterOuterAlt(localctx, 6)
		{
			p.SetState(201)
			p.RangeOperator()
		}

	case 7:
		p.EnterOuterAlt(localctx, 7)
		{
			p.SetState(202)
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

	p.SetState(209)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case FqlParserLimit:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(205)
			p.LimitClause()
		}

	case FqlParserSort:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(206)
			p.SortClause()
		}

	case FqlParserFilter:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(207)
			p.FilterClause()
		}

	case FqlParserCollect:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(208)
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

	p.SetState(213)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 11, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(211)
			p.VariableDeclaration()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(212)
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

	p.SetState(217)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 12, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(215)
			p.ForExpressionStatement()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(216)
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

	p.SetState(221)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case FqlParserReturn:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(219)
			p.ReturnExpression()
		}

	case FqlParserFor:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(220)
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
		p.SetState(223)
		p.Match(FqlParserFilter)
	}
	{
		p.SetState(224)
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
		p.SetState(226)
		p.Match(FqlParserLimit)
	}
	{
		p.SetState(227)
		p.LimitClauseValue()
	}
	p.SetState(230)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == FqlParserComma {
		{
			p.SetState(228)
			p.Match(FqlParserComma)
		}
		{
			p.SetState(229)
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

	p.SetState(234)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case FqlParserIntegerLiteral:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(232)
			p.Match(FqlParserIntegerLiteral)
		}

	case FqlParserParam:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(233)
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
		p.SetState(236)
		p.Match(FqlParserSort)
	}
	{
		p.SetState(237)
		p.SortClauseExpression()
	}
	p.SetState(242)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == FqlParserComma {
		{
			p.SetState(238)
			p.Match(FqlParserComma)
		}
		{
			p.SetState(239)
			p.SortClauseExpression()
		}

		p.SetState(244)
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
		p.SetState(245)
		p.expression(0)
	}
	p.SetState(247)
	p.GetErrorHandler().Sync(p)

	if p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 17, p.GetParserRuleContext()) == 1 {
		{
			p.SetState(246)
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

	p.SetState(267)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 18, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(249)
			p.Match(FqlParserCollect)
		}
		{
			p.SetState(250)
			p.CollectCounter()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(251)
			p.Match(FqlParserCollect)
		}
		{
			p.SetState(252)
			p.CollectAggregator()
		}

	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(253)
			p.Match(FqlParserCollect)
		}
		{
			p.SetState(254)
			p.CollectGrouping()
		}
		{
			p.SetState(255)
			p.CollectAggregator()
		}

	case 4:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(257)
			p.Match(FqlParserCollect)
		}
		{
			p.SetState(258)
			p.CollectGrouping()
		}
		{
			p.SetState(259)
			p.CollectGroupVariable()
		}

	case 5:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(261)
			p.Match(FqlParserCollect)
		}
		{
			p.SetState(262)
			p.CollectGrouping()
		}
		{
			p.SetState(263)
			p.CollectCounter()
		}

	case 6:
		p.EnterOuterAlt(localctx, 6)
		{
			p.SetState(265)
			p.Match(FqlParserCollect)
		}
		{
			p.SetState(266)
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
		p.SetState(269)
		p.Match(FqlParserIdentifier)
	}
	{
		p.SetState(270)
		p.Match(FqlParserAssign)
	}
	{
		p.SetState(271)
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
		p.SetState(273)
		p.CollectSelector()
	}
	p.SetState(278)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == FqlParserComma {
		{
			p.SetState(274)
			p.Match(FqlParserComma)
		}
		{
			p.SetState(275)
			p.CollectSelector()
		}

		p.SetState(280)
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
		p.SetState(281)
		p.Match(FqlParserAggregate)
	}
	{
		p.SetState(282)
		p.CollectAggregateSelector()
	}
	p.SetState(287)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == FqlParserComma {
		{
			p.SetState(283)
			p.Match(FqlParserComma)
		}
		{
			p.SetState(284)
			p.CollectAggregateSelector()
		}

		p.SetState(289)
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
		p.SetState(290)
		p.Match(FqlParserIdentifier)
	}
	{
		p.SetState(291)
		p.Match(FqlParserAssign)
	}
	{
		p.SetState(292)
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

	p.SetState(302)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 22, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(294)
			p.Match(FqlParserInto)
		}
		{
			p.SetState(295)
			p.CollectSelector()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(296)
			p.Match(FqlParserInto)
		}
		{
			p.SetState(297)
			p.Match(FqlParserIdentifier)
		}
		p.SetState(300)
		p.GetErrorHandler().Sync(p)

		if p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 21, p.GetParserRuleContext()) == 1 {
			{
				p.SetState(298)
				p.Match(FqlParserKeep)
			}
			{
				p.SetState(299)
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
		p.SetState(304)
		p.Match(FqlParserWith)
	}
	{
		p.SetState(305)
		p.Match(FqlParserCount)
	}
	{
		p.SetState(306)
		p.Match(FqlParserInto)
	}
	{
		p.SetState(307)
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
	p.EnterRule(localctx, 56, FqlParserRULE_variableDeclaration)

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

	p.SetState(324)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 23, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(309)
			p.Match(FqlParserLet)
		}
		{
			p.SetState(310)
			p.Match(FqlParserIdentifier)
		}
		{
			p.SetState(311)
			p.Match(FqlParserAssign)
		}
		{
			p.SetState(312)
			p.expression(0)
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(313)
			p.Match(FqlParserLet)
		}
		{
			p.SetState(314)
			p.Match(FqlParserIdentifier)
		}
		{
			p.SetState(315)
			p.Match(FqlParserAssign)
		}
		{
			p.SetState(316)
			p.Match(FqlParserOpenParen)
		}
		{
			p.SetState(317)
			p.ForExpression()
		}
		{
			p.SetState(318)
			p.Match(FqlParserCloseParen)
		}

	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(320)
			p.Match(FqlParserLet)
		}
		{
			p.SetState(321)
			p.Match(FqlParserIdentifier)
		}
		{
			p.SetState(322)
			p.Match(FqlParserAssign)
		}
		{
			p.SetState(323)
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
	p.EnterRule(localctx, 58, FqlParserRULE_param)

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
		p.SetState(326)
		p.Match(FqlParserParam)
	}
	{
		p.SetState(327)
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
	p.EnterRule(localctx, 60, FqlParserRULE_variable)

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
	p.EnterRule(localctx, 62, FqlParserRULE_rangeOperator)

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
	p.SetState(334)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case FqlParserIntegerLiteral:
		{
			p.SetState(331)
			p.IntegerLiteral()
		}

	case FqlParserIdentifier:
		{
			p.SetState(332)
			p.Variable()
		}

	case FqlParserParam:
		{
			p.SetState(333)
			p.Param()
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}
	{
		p.SetState(336)
		p.Match(FqlParserRange)
	}
	p.SetState(340)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case FqlParserIntegerLiteral:
		{
			p.SetState(337)
			p.IntegerLiteral()
		}

	case FqlParserIdentifier:
		{
			p.SetState(338)
			p.Variable()
		}

	case FqlParserParam:
		{
			p.SetState(339)
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
	p.EnterRule(localctx, 64, FqlParserRULE_arrayLiteral)
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
		p.SetState(342)
		p.Match(FqlParserOpenBracket)
	}
	p.SetState(344)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if (((_la)&-(0x1f+1)) == 0 && ((1<<uint(_la))&((1<<FqlParserOpenBracket)|(1<<FqlParserOpenParen)|(1<<FqlParserOpenBrace)|(1<<FqlParserPlus)|(1<<FqlParserMinus)|(1<<FqlParserAnd)|(1<<FqlParserOr))) != 0) || (((_la-35)&-(0x1f+1)) == 0 && ((1<<uint((_la-35)))&((1<<(FqlParserFor-35))|(1<<(FqlParserReturn-35))|(1<<(FqlParserDistinct-35))|(1<<(FqlParserFilter-35))|(1<<(FqlParserSort-35))|(1<<(FqlParserLimit-35))|(1<<(FqlParserLet-35))|(1<<(FqlParserCollect-35))|(1<<(FqlParserSortDirection-35))|(1<<(FqlParserNone-35))|(1<<(FqlParserNull-35))|(1<<(FqlParserBooleanLiteral-35))|(1<<(FqlParserUse-35))|(1<<(FqlParserInto-35))|(1<<(FqlParserKeep-35))|(1<<(FqlParserWith-35))|(1<<(FqlParserCount-35))|(1<<(FqlParserAll-35))|(1<<(FqlParserAny-35))|(1<<(FqlParserAggregate-35))|(1<<(FqlParserLike-35))|(1<<(FqlParserNot-35))|(1<<(FqlParserIn-35))|(1<<(FqlParserParam-35))|(1<<(FqlParserIdentifier-35))|(1<<(FqlParserStringLiteral-35))|(1<<(FqlParserIntegerLiteral-35))|(1<<(FqlParserFloatLiteral-35))|(1<<(FqlParserNamespaceSegment-35)))) != 0) {
		{
			p.SetState(343)
			p.ArrayElementList()
		}

	}
	{
		p.SetState(346)
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
	p.EnterRule(localctx, 66, FqlParserRULE_objectLiteral)
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
		p.SetState(348)
		p.Match(FqlParserOpenBrace)
	}
	p.SetState(357)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == FqlParserOpenBracket || (((_la-58)&-(0x1f+1)) == 0 && ((1<<uint((_la-58)))&((1<<(FqlParserParam-58))|(1<<(FqlParserIdentifier-58))|(1<<(FqlParserStringLiteral-58)))) != 0) {
		{
			p.SetState(349)
			p.PropertyAssignment()
		}
		p.SetState(354)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 27, p.GetParserRuleContext())

		for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
			if _alt == 1 {
				{
					p.SetState(350)
					p.Match(FqlParserComma)
				}
				{
					p.SetState(351)
					p.PropertyAssignment()
				}

			}
			p.SetState(356)
			p.GetErrorHandler().Sync(p)
			_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 27, p.GetParserRuleContext())
		}

	}
	p.SetState(360)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == FqlParserComma {
		{
			p.SetState(359)
			p.Match(FqlParserComma)
		}

	}
	{
		p.SetState(362)
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
	p.EnterRule(localctx, 68, FqlParserRULE_booleanLiteral)

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
		p.SetState(364)
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
	p.EnterRule(localctx, 70, FqlParserRULE_stringLiteral)

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
		p.SetState(366)
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
	p.EnterRule(localctx, 72, FqlParserRULE_integerLiteral)

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
		p.SetState(368)
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
	p.EnterRule(localctx, 74, FqlParserRULE_floatLiteral)

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
		p.SetState(370)
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
	p.EnterRule(localctx, 76, FqlParserRULE_noneLiteral)
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
		p.SetState(372)
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
	p.EnterRule(localctx, 78, FqlParserRULE_arrayElementList)
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
		p.SetState(374)
		p.expression(0)
	}
	p.SetState(383)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == FqlParserComma {
		p.SetState(376)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		for ok := true; ok; ok = _la == FqlParserComma {
			{
				p.SetState(375)
				p.Match(FqlParserComma)
			}

			p.SetState(378)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)
		}
		{
			p.SetState(380)
			p.expression(0)
		}

		p.SetState(385)
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
	p.EnterRule(localctx, 80, FqlParserRULE_propertyAssignment)

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

	p.SetState(395)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 32, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(386)
			p.PropertyName()
		}
		{
			p.SetState(387)
			p.Match(FqlParserColon)
		}
		{
			p.SetState(388)
			p.expression(0)
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(390)
			p.ComputedPropertyName()
		}
		{
			p.SetState(391)
			p.Match(FqlParserColon)
		}
		{
			p.SetState(392)
			p.expression(0)
		}

	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(394)
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
	p.EnterRule(localctx, 82, FqlParserRULE_shorthandPropertyName)

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
		p.SetState(397)
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
	p.EnterRule(localctx, 84, FqlParserRULE_computedPropertyName)

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
		p.SetState(399)
		p.Match(FqlParserOpenBracket)
	}
	{
		p.SetState(400)
		p.expression(0)
	}
	{
		p.SetState(401)
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
	p.EnterRule(localctx, 86, FqlParserRULE_propertyName)

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

	p.SetState(406)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case FqlParserIdentifier:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(403)
			p.Match(FqlParserIdentifier)
		}

	case FqlParserStringLiteral:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(404)
			p.StringLiteral()
		}

	case FqlParserParam:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(405)
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
	p.EnterRule(localctx, 88, FqlParserRULE_expressionGroup)

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
		p.SetState(408)
		p.Match(FqlParserOpenParen)
	}
	{
		p.SetState(409)
		p.expression(0)
	}
	{
		p.SetState(410)
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
	p.EnterRule(localctx, 90, FqlParserRULE_namespaceIdentifier)

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
		p.SetState(412)
		p.Namespace()
	}
	{
		p.SetState(413)
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
	p.EnterRule(localctx, 92, FqlParserRULE_namespace)
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
	p.SetState(418)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == FqlParserNamespaceSegment {
		{
			p.SetState(415)
			p.Match(FqlParserNamespaceSegment)
		}

		p.SetState(420)
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
	p.EnterRule(localctx, 94, FqlParserRULE_functionIdentifier)
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
		p.SetState(421)
		_la = p.GetTokenStream().LA(1)

		if !(((_la-28)&-(0x1f+1)) == 0 && ((1<<uint((_la-28)))&((1<<(FqlParserAnd-28))|(1<<(FqlParserOr-28))|(1<<(FqlParserFor-28))|(1<<(FqlParserReturn-28))|(1<<(FqlParserDistinct-28))|(1<<(FqlParserFilter-28))|(1<<(FqlParserSort-28))|(1<<(FqlParserLimit-28))|(1<<(FqlParserLet-28))|(1<<(FqlParserCollect-28))|(1<<(FqlParserSortDirection-28))|(1<<(FqlParserNone-28))|(1<<(FqlParserNull-28))|(1<<(FqlParserBooleanLiteral-28))|(1<<(FqlParserUse-28))|(1<<(FqlParserInto-28))|(1<<(FqlParserKeep-28))|(1<<(FqlParserWith-28))|(1<<(FqlParserCount-28))|(1<<(FqlParserAll-28))|(1<<(FqlParserAny-28))|(1<<(FqlParserAggregate-28))|(1<<(FqlParserLike-28))|(1<<(FqlParserNot-28))|(1<<(FqlParserIn-28))|(1<<(FqlParserIdentifier-28)))) != 0) {
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
	p.EnterRule(localctx, 96, FqlParserRULE_functionCallExpression)

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
		p.Namespace()
	}
	{
		p.SetState(424)
		p.FunctionIdentifier()
	}
	{
		p.SetState(425)
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
	p.EnterRule(localctx, 98, FqlParserRULE_member)

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

	p.SetState(430)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 35, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(427)
			p.Match(FqlParserIdentifier)
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(428)
			p.FunctionCallExpression()
		}

	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(429)
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
	p.EnterRule(localctx, 100, FqlParserRULE_memberPath)

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

	p.SetState(471)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case FqlParserDot:
		p.EnterOuterAlt(localctx, 1)
		p.SetState(440)
		p.GetErrorHandler().Sync(p)
		_alt = 1
		for ok := true; ok; ok = _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
			switch _alt {
			case 1:
				{
					p.SetState(432)
					p.Match(FqlParserDot)
				}
				{
					p.SetState(433)
					p.PropertyName()
				}
				p.SetState(437)
				p.GetErrorHandler().Sync(p)
				_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 36, p.GetParserRuleContext())

				for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
					if _alt == 1 {
						{
							p.SetState(434)
							p.ComputedPropertyName()
						}

					}
					p.SetState(439)
					p.GetErrorHandler().Sync(p)
					_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 36, p.GetParserRuleContext())
				}

			default:
				panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
			}

			p.SetState(442)
			p.GetErrorHandler().Sync(p)
			_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 37, p.GetParserRuleContext())
		}

	case FqlParserOpenBracket:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(444)
			p.ComputedPropertyName()
		}
		p.SetState(455)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 39, p.GetParserRuleContext())

		for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
			if _alt == 1 {
				{
					p.SetState(445)
					p.Match(FqlParserDot)
				}
				{
					p.SetState(446)
					p.PropertyName()
				}
				p.SetState(450)
				p.GetErrorHandler().Sync(p)
				_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 38, p.GetParserRuleContext())

				for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
					if _alt == 1 {
						{
							p.SetState(447)
							p.ComputedPropertyName()
						}

					}
					p.SetState(452)
					p.GetErrorHandler().Sync(p)
					_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 38, p.GetParserRuleContext())
				}

			}
			p.SetState(457)
			p.GetErrorHandler().Sync(p)
			_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 39, p.GetParserRuleContext())
		}
		p.SetState(468)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 41, p.GetParserRuleContext())

		for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
			if _alt == 1 {
				{
					p.SetState(458)
					p.ComputedPropertyName()
				}
				p.SetState(463)
				p.GetErrorHandler().Sync(p)
				_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 40, p.GetParserRuleContext())

				for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
					if _alt == 1 {
						{
							p.SetState(459)
							p.Match(FqlParserDot)
						}
						{
							p.SetState(460)
							p.PropertyName()
						}

					}
					p.SetState(465)
					p.GetErrorHandler().Sync(p)
					_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 40, p.GetParserRuleContext())
				}

			}
			p.SetState(470)
			p.GetErrorHandler().Sync(p)
			_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 41, p.GetParserRuleContext())
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
	p.EnterRule(localctx, 102, FqlParserRULE_memberExpression)

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
		p.SetState(473)
		p.Member()
	}
	{
		p.SetState(474)
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
	p.EnterRule(localctx, 104, FqlParserRULE_arguments)
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
		p.SetState(476)
		p.Match(FqlParserOpenParen)
	}
	p.SetState(485)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if (((_la)&-(0x1f+1)) == 0 && ((1<<uint(_la))&((1<<FqlParserOpenBracket)|(1<<FqlParserOpenParen)|(1<<FqlParserOpenBrace)|(1<<FqlParserPlus)|(1<<FqlParserMinus)|(1<<FqlParserAnd)|(1<<FqlParserOr))) != 0) || (((_la-35)&-(0x1f+1)) == 0 && ((1<<uint((_la-35)))&((1<<(FqlParserFor-35))|(1<<(FqlParserReturn-35))|(1<<(FqlParserDistinct-35))|(1<<(FqlParserFilter-35))|(1<<(FqlParserSort-35))|(1<<(FqlParserLimit-35))|(1<<(FqlParserLet-35))|(1<<(FqlParserCollect-35))|(1<<(FqlParserSortDirection-35))|(1<<(FqlParserNone-35))|(1<<(FqlParserNull-35))|(1<<(FqlParserBooleanLiteral-35))|(1<<(FqlParserUse-35))|(1<<(FqlParserInto-35))|(1<<(FqlParserKeep-35))|(1<<(FqlParserWith-35))|(1<<(FqlParserCount-35))|(1<<(FqlParserAll-35))|(1<<(FqlParserAny-35))|(1<<(FqlParserAggregate-35))|(1<<(FqlParserLike-35))|(1<<(FqlParserNot-35))|(1<<(FqlParserIn-35))|(1<<(FqlParserParam-35))|(1<<(FqlParserIdentifier-35))|(1<<(FqlParserStringLiteral-35))|(1<<(FqlParserIntegerLiteral-35))|(1<<(FqlParserFloatLiteral-35))|(1<<(FqlParserNamespaceSegment-35)))) != 0) {
		{
			p.SetState(477)
			p.expression(0)
		}
		p.SetState(482)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		for _la == FqlParserComma {
			{
				p.SetState(478)
				p.Match(FqlParserComma)
			}
			{
				p.SetState(479)
				p.expression(0)
			}

			p.SetState(484)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)
		}

	}
	{
		p.SetState(487)
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
	_startState := 106
	p.EnterRecursionRule(localctx, 106, FqlParserRULE_expression, _p)
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
	p.SetState(506)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 45, p.GetParserRuleContext()) {
	case 1:
		{
			p.SetState(490)
			p.UnaryOperator()
		}
		{
			p.SetState(491)
			p.expression(23)
		}

	case 2:
		{
			p.SetState(493)
			p.FunctionCallExpression()
		}

	case 3:
		{
			p.SetState(494)
			p.ExpressionGroup()
		}

	case 4:
		{
			p.SetState(495)
			p.RangeOperator()
		}

	case 5:
		{
			p.SetState(496)
			p.StringLiteral()
		}

	case 6:
		{
			p.SetState(497)
			p.IntegerLiteral()
		}

	case 7:
		{
			p.SetState(498)
			p.FloatLiteral()
		}

	case 8:
		{
			p.SetState(499)
			p.BooleanLiteral()
		}

	case 9:
		{
			p.SetState(500)
			p.ArrayLiteral()
		}

	case 10:
		{
			p.SetState(501)
			p.ObjectLiteral()
		}

	case 11:
		{
			p.SetState(502)
			p.Variable()
		}

	case 12:
		{
			p.SetState(503)
			p.MemberExpression()
		}

	case 13:
		{
			p.SetState(504)
			p.NoneLiteral()
		}

	case 14:
		{
			p.SetState(505)
			p.Param()
		}

	}
	p.GetParserRuleContext().SetStop(p.GetTokenStream().LT(-1))
	p.SetState(553)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 49, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			if p.GetParseListeners() != nil {
				p.TriggerExitRuleEvent()
			}
			_prevctx = localctx
			p.SetState(551)
			p.GetErrorHandler().Sync(p)
			switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 48, p.GetParserRuleContext()) {
			case 1:
				localctx = NewExpressionContext(p, _parentctx, _parentState)
				p.PushNewRecursionContext(localctx, _startState, FqlParserRULE_expression)
				p.SetState(508)

				if !(p.Precpred(p.GetParserRuleContext(), 22)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 22)", ""))
				}
				{
					p.SetState(509)
					p.MultiplicativeOperator()
				}
				{
					p.SetState(510)
					p.expression(23)
				}

			case 2:
				localctx = NewExpressionContext(p, _parentctx, _parentState)
				p.PushNewRecursionContext(localctx, _startState, FqlParserRULE_expression)
				p.SetState(512)

				if !(p.Precpred(p.GetParserRuleContext(), 21)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 21)", ""))
				}
				{
					p.SetState(513)
					p.AdditiveOperator()
				}
				{
					p.SetState(514)
					p.expression(22)
				}

			case 3:
				localctx = NewExpressionContext(p, _parentctx, _parentState)
				p.PushNewRecursionContext(localctx, _startState, FqlParserRULE_expression)
				p.SetState(516)

				if !(p.Precpred(p.GetParserRuleContext(), 18)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 18)", ""))
				}
				{
					p.SetState(517)
					p.ArrayOperator()
				}
				p.SetState(520)
				p.GetErrorHandler().Sync(p)

				switch p.GetTokenStream().LA(1) {
				case FqlParserNot, FqlParserIn:
					{
						p.SetState(518)
						p.InOperator()
					}

				case FqlParserGt, FqlParserLt, FqlParserEq, FqlParserGte, FqlParserLte, FqlParserNeq:
					{
						p.SetState(519)
						p.EqualityOperator()
					}

				default:
					panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
				}
				{
					p.SetState(522)
					p.expression(19)
				}

			case 4:
				localctx = NewExpressionContext(p, _parentctx, _parentState)
				p.PushNewRecursionContext(localctx, _startState, FqlParserRULE_expression)
				p.SetState(524)

				if !(p.Precpred(p.GetParserRuleContext(), 17)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 17)", ""))
				}
				{
					p.SetState(525)
					p.InOperator()
				}
				{
					p.SetState(526)
					p.expression(18)
				}

			case 5:
				localctx = NewExpressionContext(p, _parentctx, _parentState)
				p.PushNewRecursionContext(localctx, _startState, FqlParserRULE_expression)
				p.SetState(528)

				if !(p.Precpred(p.GetParserRuleContext(), 16)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 16)", ""))
				}
				{
					p.SetState(529)
					p.EqualityOperator()
				}
				{
					p.SetState(530)
					p.expression(17)
				}

			case 6:
				localctx = NewExpressionContext(p, _parentctx, _parentState)
				p.PushNewRecursionContext(localctx, _startState, FqlParserRULE_expression)
				p.SetState(532)

				if !(p.Precpred(p.GetParserRuleContext(), 15)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 15)", ""))
				}
				{
					p.SetState(533)
					p.RegexpOperator()
				}
				{
					p.SetState(534)
					p.expression(16)
				}

			case 7:
				localctx = NewExpressionContext(p, _parentctx, _parentState)
				p.PushNewRecursionContext(localctx, _startState, FqlParserRULE_expression)
				p.SetState(536)

				if !(p.Precpred(p.GetParserRuleContext(), 14)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 14)", ""))
				}
				{
					p.SetState(537)
					p.LogicalAndOperator()
				}
				{
					p.SetState(538)
					p.expression(15)
				}

			case 8:
				localctx = NewExpressionContext(p, _parentctx, _parentState)
				p.PushNewRecursionContext(localctx, _startState, FqlParserRULE_expression)
				p.SetState(540)

				if !(p.Precpred(p.GetParserRuleContext(), 13)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 13)", ""))
				}
				{
					p.SetState(541)
					p.LogicalOrOperator()
				}
				{
					p.SetState(542)
					p.expression(14)
				}

			case 9:
				localctx = NewExpressionContext(p, _parentctx, _parentState)
				p.PushNewRecursionContext(localctx, _startState, FqlParserRULE_expression)
				p.SetState(544)

				if !(p.Precpred(p.GetParserRuleContext(), 12)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 12)", ""))
				}
				{
					p.SetState(545)
					p.Match(FqlParserQuestionMark)
				}
				p.SetState(547)
				p.GetErrorHandler().Sync(p)
				_la = p.GetTokenStream().LA(1)

				if (((_la)&-(0x1f+1)) == 0 && ((1<<uint(_la))&((1<<FqlParserOpenBracket)|(1<<FqlParserOpenParen)|(1<<FqlParserOpenBrace)|(1<<FqlParserPlus)|(1<<FqlParserMinus)|(1<<FqlParserAnd)|(1<<FqlParserOr))) != 0) || (((_la-35)&-(0x1f+1)) == 0 && ((1<<uint((_la-35)))&((1<<(FqlParserFor-35))|(1<<(FqlParserReturn-35))|(1<<(FqlParserDistinct-35))|(1<<(FqlParserFilter-35))|(1<<(FqlParserSort-35))|(1<<(FqlParserLimit-35))|(1<<(FqlParserLet-35))|(1<<(FqlParserCollect-35))|(1<<(FqlParserSortDirection-35))|(1<<(FqlParserNone-35))|(1<<(FqlParserNull-35))|(1<<(FqlParserBooleanLiteral-35))|(1<<(FqlParserUse-35))|(1<<(FqlParserInto-35))|(1<<(FqlParserKeep-35))|(1<<(FqlParserWith-35))|(1<<(FqlParserCount-35))|(1<<(FqlParserAll-35))|(1<<(FqlParserAny-35))|(1<<(FqlParserAggregate-35))|(1<<(FqlParserLike-35))|(1<<(FqlParserNot-35))|(1<<(FqlParserIn-35))|(1<<(FqlParserParam-35))|(1<<(FqlParserIdentifier-35))|(1<<(FqlParserStringLiteral-35))|(1<<(FqlParserIntegerLiteral-35))|(1<<(FqlParserFloatLiteral-35))|(1<<(FqlParserNamespaceSegment-35)))) != 0) {
					{
						p.SetState(546)
						p.expression(0)
					}

				}
				{
					p.SetState(549)
					p.Match(FqlParserColon)
				}
				{
					p.SetState(550)
					p.expression(13)
				}

			}

		}
		p.SetState(555)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 49, p.GetParserRuleContext())
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
	p.EnterRule(localctx, 108, FqlParserRULE_forTernaryExpression)
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

	p.SetState(584)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 51, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(556)
			p.expression(0)
		}
		{
			p.SetState(557)
			p.Match(FqlParserQuestionMark)
		}
		p.SetState(559)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		if (((_la)&-(0x1f+1)) == 0 && ((1<<uint(_la))&((1<<FqlParserOpenBracket)|(1<<FqlParserOpenParen)|(1<<FqlParserOpenBrace)|(1<<FqlParserPlus)|(1<<FqlParserMinus)|(1<<FqlParserAnd)|(1<<FqlParserOr))) != 0) || (((_la-35)&-(0x1f+1)) == 0 && ((1<<uint((_la-35)))&((1<<(FqlParserFor-35))|(1<<(FqlParserReturn-35))|(1<<(FqlParserDistinct-35))|(1<<(FqlParserFilter-35))|(1<<(FqlParserSort-35))|(1<<(FqlParserLimit-35))|(1<<(FqlParserLet-35))|(1<<(FqlParserCollect-35))|(1<<(FqlParserSortDirection-35))|(1<<(FqlParserNone-35))|(1<<(FqlParserNull-35))|(1<<(FqlParserBooleanLiteral-35))|(1<<(FqlParserUse-35))|(1<<(FqlParserInto-35))|(1<<(FqlParserKeep-35))|(1<<(FqlParserWith-35))|(1<<(FqlParserCount-35))|(1<<(FqlParserAll-35))|(1<<(FqlParserAny-35))|(1<<(FqlParserAggregate-35))|(1<<(FqlParserLike-35))|(1<<(FqlParserNot-35))|(1<<(FqlParserIn-35))|(1<<(FqlParserParam-35))|(1<<(FqlParserIdentifier-35))|(1<<(FqlParserStringLiteral-35))|(1<<(FqlParserIntegerLiteral-35))|(1<<(FqlParserFloatLiteral-35))|(1<<(FqlParserNamespaceSegment-35)))) != 0) {
			{
				p.SetState(558)
				p.expression(0)
			}

		}
		{
			p.SetState(561)
			p.Match(FqlParserColon)
		}
		{
			p.SetState(562)
			p.Match(FqlParserOpenParen)
		}
		{
			p.SetState(563)
			p.ForExpression()
		}
		{
			p.SetState(564)
			p.Match(FqlParserCloseParen)
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(566)
			p.expression(0)
		}
		{
			p.SetState(567)
			p.Match(FqlParserQuestionMark)
		}
		{
			p.SetState(568)
			p.Match(FqlParserOpenParen)
		}
		{
			p.SetState(569)
			p.ForExpression()
		}
		{
			p.SetState(570)
			p.Match(FqlParserCloseParen)
		}
		{
			p.SetState(571)
			p.Match(FqlParserColon)
		}
		{
			p.SetState(572)
			p.expression(0)
		}

	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(574)
			p.expression(0)
		}
		{
			p.SetState(575)
			p.Match(FqlParserQuestionMark)
		}
		{
			p.SetState(576)
			p.Match(FqlParserOpenParen)
		}
		{
			p.SetState(577)
			p.ForExpression()
		}
		{
			p.SetState(578)
			p.Match(FqlParserCloseParen)
		}
		{
			p.SetState(579)
			p.Match(FqlParserColon)
		}
		{
			p.SetState(580)
			p.Match(FqlParserOpenParen)
		}
		{
			p.SetState(581)
			p.ForExpression()
		}
		{
			p.SetState(582)
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
	p.EnterRule(localctx, 110, FqlParserRULE_arrayOperator)
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
		p.SetState(586)
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
	p.EnterRule(localctx, 112, FqlParserRULE_inOperator)

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

	p.SetState(591)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case FqlParserIn:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(588)
			p.Match(FqlParserIn)
		}

	case FqlParserNot:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(589)
			p.Match(FqlParserNot)
		}
		{
			p.SetState(590)
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
	p.EnterRule(localctx, 114, FqlParserRULE_equalityOperator)
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
		p.SetState(593)
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
	p.EnterRule(localctx, 116, FqlParserRULE_regexpOperator)
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
		p.SetState(595)
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
	p.EnterRule(localctx, 118, FqlParserRULE_logicalAndOperator)

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
		p.SetState(597)
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
	p.EnterRule(localctx, 120, FqlParserRULE_logicalOrOperator)

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
		p.SetState(599)
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
	p.EnterRule(localctx, 122, FqlParserRULE_multiplicativeOperator)
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
		p.SetState(601)
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
	p.EnterRule(localctx, 124, FqlParserRULE_additiveOperator)
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
		p.SetState(603)
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
	p.EnterRule(localctx, 126, FqlParserRULE_unaryOperator)
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
		p.SetState(605)
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
	case 53:
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
