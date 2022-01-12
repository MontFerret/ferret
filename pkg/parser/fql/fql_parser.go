// Code generated from antlr/FqlParser.g4 by ANTLR 4.9.3. DO NOT EDIT.

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
	3, 24715, 42794, 33075, 47597, 16764, 15335, 30598, 22884, 3, 74, 649,
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
	4, 66, 9, 66, 4, 67, 9, 67, 4, 68, 9, 68, 4, 69, 9, 69, 4, 70, 9, 70, 4,
	71, 9, 71, 4, 72, 9, 72, 4, 73, 9, 73, 3, 2, 7, 2, 148, 10, 2, 12, 2, 14,
	2, 151, 11, 2, 3, 2, 3, 2, 3, 3, 3, 3, 3, 4, 3, 4, 3, 5, 3, 5, 3, 5, 3,
	6, 7, 6, 163, 10, 6, 12, 6, 14, 6, 166, 11, 6, 3, 6, 3, 6, 3, 7, 3, 7,
	3, 7, 5, 7, 173, 10, 7, 3, 8, 3, 8, 5, 8, 177, 10, 8, 3, 9, 3, 9, 3, 9,
	3, 9, 3, 9, 3, 9, 3, 9, 3, 9, 3, 9, 5, 9, 188, 10, 9, 3, 10, 3, 10, 5,
	10, 192, 10, 10, 3, 10, 3, 10, 3, 11, 3, 11, 3, 11, 3, 11, 5, 11, 200,
	10, 11, 3, 11, 3, 11, 3, 11, 7, 11, 205, 10, 11, 12, 11, 14, 11, 208, 11,
	11, 3, 11, 3, 11, 3, 11, 3, 11, 3, 11, 5, 11, 215, 10, 11, 3, 11, 3, 11,
	3, 11, 7, 11, 220, 10, 11, 12, 11, 14, 11, 223, 11, 11, 3, 11, 3, 11, 5,
	11, 227, 10, 11, 3, 12, 3, 12, 3, 12, 3, 12, 3, 12, 3, 12, 3, 12, 5, 12,
	236, 10, 12, 3, 13, 3, 13, 3, 13, 3, 13, 5, 13, 242, 10, 13, 3, 14, 3,
	14, 5, 14, 246, 10, 14, 3, 15, 3, 15, 5, 15, 250, 10, 15, 3, 16, 3, 16,
	5, 16, 254, 10, 16, 3, 17, 3, 17, 3, 17, 3, 18, 3, 18, 3, 18, 3, 18, 5,
	18, 263, 10, 18, 3, 19, 3, 19, 3, 19, 3, 19, 3, 19, 5, 19, 270, 10, 19,
	3, 20, 3, 20, 3, 20, 3, 20, 7, 20, 276, 10, 20, 12, 20, 14, 20, 279, 11,
	20, 3, 21, 3, 21, 5, 21, 283, 10, 21, 3, 22, 3, 22, 3, 22, 3, 22, 3, 22,
	3, 22, 3, 22, 3, 22, 3, 22, 3, 22, 3, 22, 3, 22, 3, 22, 3, 22, 3, 22, 3,
	22, 3, 22, 3, 22, 5, 22, 303, 10, 22, 3, 23, 3, 23, 3, 23, 3, 23, 3, 24,
	3, 24, 3, 24, 7, 24, 312, 10, 24, 12, 24, 14, 24, 315, 11, 24, 3, 25, 3,
	25, 3, 25, 3, 25, 7, 25, 321, 10, 25, 12, 25, 14, 25, 324, 11, 25, 3, 26,
	3, 26, 3, 26, 3, 26, 3, 27, 3, 27, 3, 27, 3, 27, 3, 27, 3, 27, 5, 27, 336,
	10, 27, 5, 27, 338, 10, 27, 3, 28, 3, 28, 3, 28, 3, 28, 3, 28, 3, 29, 3,
	29, 3, 29, 3, 29, 3, 29, 3, 29, 5, 29, 351, 10, 29, 3, 29, 5, 29, 354,
	10, 29, 3, 29, 5, 29, 357, 10, 29, 3, 30, 3, 30, 3, 30, 3, 30, 3, 30, 5,
	30, 364, 10, 30, 3, 31, 3, 31, 3, 31, 5, 31, 369, 10, 31, 3, 32, 3, 32,
	3, 32, 3, 33, 3, 33, 3, 33, 3, 33, 3, 33, 3, 33, 5, 33, 380, 10, 33, 3,
	34, 3, 34, 3, 34, 3, 34, 5, 34, 386, 10, 34, 3, 35, 3, 35, 5, 35, 390,
	10, 35, 3, 36, 3, 36, 3, 36, 3, 36, 3, 36, 3, 36, 3, 36, 5, 36, 399, 10,
	36, 3, 37, 3, 37, 5, 37, 403, 10, 37, 3, 37, 3, 37, 3, 38, 3, 38, 3, 38,
	3, 38, 7, 38, 411, 10, 38, 12, 38, 14, 38, 414, 11, 38, 3, 38, 5, 38, 417,
	10, 38, 5, 38, 419, 10, 38, 3, 38, 3, 38, 3, 39, 3, 39, 3, 40, 3, 40, 3,
	41, 3, 41, 3, 42, 3, 42, 3, 43, 3, 43, 3, 44, 3, 44, 3, 44, 3, 44, 3, 44,
	3, 44, 3, 44, 3, 44, 3, 44, 5, 44, 442, 10, 44, 3, 45, 3, 45, 3, 45, 3,
	45, 3, 46, 3, 46, 3, 46, 3, 46, 3, 46, 5, 46, 453, 10, 46, 3, 47, 3, 47,
	3, 47, 3, 48, 7, 48, 459, 10, 48, 12, 48, 14, 48, 462, 11, 48, 3, 49, 3,
	49, 6, 49, 466, 10, 49, 13, 49, 14, 49, 467, 3, 50, 3, 50, 3, 50, 3, 50,
	3, 50, 5, 50, 475, 10, 50, 3, 51, 3, 51, 5, 51, 479, 10, 51, 3, 52, 3,
	52, 3, 52, 3, 52, 5, 52, 485, 10, 52, 3, 52, 3, 52, 3, 53, 3, 53, 3, 53,
	5, 53, 492, 10, 53, 3, 54, 3, 54, 3, 54, 7, 54, 497, 10, 54, 12, 54, 14,
	54, 500, 11, 54, 3, 54, 5, 54, 503, 10, 54, 3, 55, 5, 55, 506, 10, 55,
	3, 55, 3, 55, 3, 55, 3, 55, 3, 55, 5, 55, 513, 10, 55, 3, 55, 5, 55, 516,
	10, 55, 3, 56, 3, 56, 3, 57, 3, 57, 3, 58, 3, 58, 3, 58, 3, 58, 3, 59,
	3, 59, 3, 59, 5, 59, 529, 10, 59, 3, 60, 3, 60, 3, 60, 3, 60, 3, 60, 5,
	60, 536, 10, 60, 3, 60, 3, 60, 3, 60, 3, 60, 3, 60, 3, 60, 3, 60, 3, 60,
	3, 60, 3, 60, 3, 60, 5, 60, 549, 10, 60, 3, 60, 3, 60, 7, 60, 553, 10,
	60, 12, 60, 14, 60, 556, 11, 60, 3, 61, 3, 61, 3, 61, 3, 61, 3, 61, 3,
	61, 3, 61, 3, 61, 3, 61, 3, 61, 3, 61, 3, 61, 3, 61, 3, 61, 3, 61, 3, 61,
	3, 61, 3, 61, 3, 61, 7, 61, 577, 10, 61, 12, 61, 14, 61, 580, 11, 61, 3,
	62, 3, 62, 3, 62, 3, 62, 3, 62, 3, 62, 3, 62, 3, 62, 3, 62, 3, 62, 3, 62,
	5, 62, 593, 10, 62, 3, 62, 3, 62, 5, 62, 597, 10, 62, 5, 62, 599, 10, 62,
	3, 62, 3, 62, 3, 62, 3, 62, 3, 62, 3, 62, 3, 62, 3, 62, 3, 62, 3, 62, 3,
	62, 3, 62, 7, 62, 613, 10, 62, 12, 62, 14, 62, 616, 11, 62, 3, 63, 3, 63,
	3, 63, 5, 63, 621, 10, 63, 3, 64, 3, 64, 3, 65, 5, 65, 626, 10, 65, 3,
	65, 3, 65, 3, 66, 5, 66, 631, 10, 66, 3, 66, 3, 66, 3, 67, 3, 67, 3, 68,
	3, 68, 3, 69, 3, 69, 3, 70, 3, 70, 3, 71, 3, 71, 3, 72, 3, 72, 3, 73, 3,
	73, 3, 73, 2, 5, 118, 120, 122, 74, 2, 4, 6, 8, 10, 12, 14, 16, 18, 20,
	22, 24, 26, 28, 30, 32, 34, 36, 38, 40, 42, 44, 46, 48, 50, 52, 54, 56,
	58, 60, 62, 64, 66, 68, 70, 72, 74, 76, 78, 80, 82, 84, 86, 88, 90, 92,
	94, 96, 98, 100, 102, 104, 106, 108, 110, 112, 114, 116, 118, 120, 122,
	124, 126, 128, 130, 132, 134, 136, 138, 140, 142, 144, 2, 12, 3, 2, 68,
	69, 3, 2, 50, 51, 6, 2, 30, 31, 40, 46, 48, 49, 54, 61, 6, 2, 37, 39, 47,
	47, 50, 53, 62, 66, 4, 2, 50, 50, 58, 59, 3, 2, 17, 22, 4, 2, 26, 27, 63,
	63, 3, 2, 35, 36, 3, 2, 23, 25, 3, 2, 26, 27, 2, 688, 2, 149, 3, 2, 2,
	2, 4, 154, 3, 2, 2, 2, 6, 156, 3, 2, 2, 2, 8, 158, 3, 2, 2, 2, 10, 164,
	3, 2, 2, 2, 12, 172, 3, 2, 2, 2, 14, 176, 3, 2, 2, 2, 16, 187, 3, 2, 2,
	2, 18, 189, 3, 2, 2, 2, 20, 226, 3, 2, 2, 2, 22, 235, 3, 2, 2, 2, 24, 241,
	3, 2, 2, 2, 26, 245, 3, 2, 2, 2, 28, 249, 3, 2, 2, 2, 30, 253, 3, 2, 2,
	2, 32, 255, 3, 2, 2, 2, 34, 258, 3, 2, 2, 2, 36, 269, 3, 2, 2, 2, 38, 271,
	3, 2, 2, 2, 40, 280, 3, 2, 2, 2, 42, 302, 3, 2, 2, 2, 44, 304, 3, 2, 2,
	2, 46, 308, 3, 2, 2, 2, 48, 316, 3, 2, 2, 2, 50, 325, 3, 2, 2, 2, 52, 337,
	3, 2, 2, 2, 54, 339, 3, 2, 2, 2, 56, 344, 3, 2, 2, 2, 58, 363, 3, 2, 2,
	2, 60, 368, 3, 2, 2, 2, 62, 370, 3, 2, 2, 2, 64, 373, 3, 2, 2, 2, 66, 385,
	3, 2, 2, 2, 68, 389, 3, 2, 2, 2, 70, 398, 3, 2, 2, 2, 72, 400, 3, 2, 2,
	2, 74, 406, 3, 2, 2, 2, 76, 422, 3, 2, 2, 2, 78, 424, 3, 2, 2, 2, 80, 426,
	3, 2, 2, 2, 82, 428, 3, 2, 2, 2, 84, 430, 3, 2, 2, 2, 86, 441, 3, 2, 2,
	2, 88, 443, 3, 2, 2, 2, 90, 452, 3, 2, 2, 2, 92, 454, 3, 2, 2, 2, 94, 460,
	3, 2, 2, 2, 96, 463, 3, 2, 2, 2, 98, 474, 3, 2, 2, 2, 100, 476, 3, 2, 2,
	2, 102, 480, 3, 2, 2, 2, 104, 491, 3, 2, 2, 2, 106, 493, 3, 2, 2, 2, 108,
	515, 3, 2, 2, 2, 110, 517, 3, 2, 2, 2, 112, 519, 3, 2, 2, 2, 114, 521,
	3, 2, 2, 2, 116, 528, 3, 2, 2, 2, 118, 535, 3, 2, 2, 2, 120, 557, 3, 2,
	2, 2, 122, 598, 3, 2, 2, 2, 124, 617, 3, 2, 2, 2, 126, 622, 3, 2, 2, 2,
	128, 625, 3, 2, 2, 2, 130, 630, 3, 2, 2, 2, 132, 634, 3, 2, 2, 2, 134,
	636, 3, 2, 2, 2, 136, 638, 3, 2, 2, 2, 138, 640, 3, 2, 2, 2, 140, 642,
	3, 2, 2, 2, 142, 644, 3, 2, 2, 2, 144, 646, 3, 2, 2, 2, 146, 148, 5, 4,
	3, 2, 147, 146, 3, 2, 2, 2, 148, 151, 3, 2, 2, 2, 149, 147, 3, 2, 2, 2,
	149, 150, 3, 2, 2, 2, 150, 152, 3, 2, 2, 2, 151, 149, 3, 2, 2, 2, 152,
	153, 5, 10, 6, 2, 153, 3, 3, 2, 2, 2, 154, 155, 5, 6, 4, 2, 155, 5, 3,
	2, 2, 2, 156, 157, 5, 8, 5, 2, 157, 7, 3, 2, 2, 2, 158, 159, 7, 53, 2,
	2, 159, 160, 5, 92, 47, 2, 160, 9, 3, 2, 2, 2, 161, 163, 5, 12, 7, 2, 162,
	161, 3, 2, 2, 2, 163, 166, 3, 2, 2, 2, 164, 162, 3, 2, 2, 2, 164, 165,
	3, 2, 2, 2, 165, 167, 3, 2, 2, 2, 166, 164, 3, 2, 2, 2, 167, 168, 5, 14,
	8, 2, 168, 11, 3, 2, 2, 2, 169, 173, 5, 16, 9, 2, 170, 173, 5, 100, 51,
	2, 171, 173, 5, 56, 29, 2, 172, 169, 3, 2, 2, 2, 172, 170, 3, 2, 2, 2,
	172, 171, 3, 2, 2, 2, 173, 13, 3, 2, 2, 2, 174, 177, 5, 18, 10, 2, 175,
	177, 5, 20, 11, 2, 176, 174, 3, 2, 2, 2, 176, 175, 3, 2, 2, 2, 177, 15,
	3, 2, 2, 2, 178, 179, 7, 47, 2, 2, 179, 180, 9, 2, 2, 2, 180, 181, 7, 33,
	2, 2, 181, 188, 5, 118, 60, 2, 182, 183, 7, 47, 2, 2, 183, 184, 5, 110,
	56, 2, 184, 185, 7, 33, 2, 2, 185, 186, 5, 118, 60, 2, 186, 188, 3, 2,
	2, 2, 187, 178, 3, 2, 2, 2, 187, 182, 3, 2, 2, 2, 188, 17, 3, 2, 2, 2,
	189, 191, 7, 38, 2, 2, 190, 192, 7, 42, 2, 2, 191, 190, 3, 2, 2, 2, 191,
	192, 3, 2, 2, 2, 192, 193, 3, 2, 2, 2, 193, 194, 5, 118, 60, 2, 194, 19,
	3, 2, 2, 2, 195, 196, 7, 37, 2, 2, 196, 199, 9, 2, 2, 2, 197, 198, 7, 10,
	2, 2, 198, 200, 7, 68, 2, 2, 199, 197, 3, 2, 2, 2, 199, 200, 3, 2, 2, 2,
	200, 201, 3, 2, 2, 2, 201, 202, 7, 64, 2, 2, 202, 206, 5, 22, 12, 2, 203,
	205, 5, 28, 15, 2, 204, 203, 3, 2, 2, 2, 205, 208, 3, 2, 2, 2, 206, 204,
	3, 2, 2, 2, 206, 207, 3, 2, 2, 2, 207, 209, 3, 2, 2, 2, 208, 206, 3, 2,
	2, 2, 209, 210, 5, 30, 16, 2, 210, 227, 3, 2, 2, 2, 211, 212, 7, 37, 2,
	2, 212, 214, 9, 2, 2, 2, 213, 215, 7, 65, 2, 2, 214, 213, 3, 2, 2, 2, 214,
	215, 3, 2, 2, 2, 215, 216, 3, 2, 2, 2, 216, 217, 7, 66, 2, 2, 217, 221,
	5, 118, 60, 2, 218, 220, 5, 28, 15, 2, 219, 218, 3, 2, 2, 2, 220, 223,
	3, 2, 2, 2, 221, 219, 3, 2, 2, 2, 221, 222, 3, 2, 2, 2, 222, 224, 3, 2,
	2, 2, 223, 221, 3, 2, 2, 2, 224, 225, 5, 30, 16, 2, 225, 227, 3, 2, 2,
	2, 226, 195, 3, 2, 2, 2, 226, 211, 3, 2, 2, 2, 227, 21, 3, 2, 2, 2, 228,
	236, 5, 100, 51, 2, 229, 236, 5, 72, 37, 2, 230, 236, 5, 74, 38, 2, 231,
	236, 5, 68, 35, 2, 232, 236, 5, 96, 49, 2, 233, 236, 5, 114, 58, 2, 234,
	236, 5, 66, 34, 2, 235, 228, 3, 2, 2, 2, 235, 229, 3, 2, 2, 2, 235, 230,
	3, 2, 2, 2, 235, 231, 3, 2, 2, 2, 235, 232, 3, 2, 2, 2, 235, 233, 3, 2,
	2, 2, 235, 234, 3, 2, 2, 2, 236, 23, 3, 2, 2, 2, 237, 242, 5, 34, 18, 2,
	238, 242, 5, 38, 20, 2, 239, 242, 5, 32, 17, 2, 240, 242, 5, 42, 22, 2,
	241, 237, 3, 2, 2, 2, 241, 238, 3, 2, 2, 2, 241, 239, 3, 2, 2, 2, 241,
	240, 3, 2, 2, 2, 242, 25, 3, 2, 2, 2, 243, 246, 5, 16, 9, 2, 244, 246,
	5, 100, 51, 2, 245, 243, 3, 2, 2, 2, 245, 244, 3, 2, 2, 2, 246, 27, 3,
	2, 2, 2, 247, 250, 5, 26, 14, 2, 248, 250, 5, 24, 13, 2, 249, 247, 3, 2,
	2, 2, 249, 248, 3, 2, 2, 2, 250, 29, 3, 2, 2, 2, 251, 254, 5, 18, 10, 2,
	252, 254, 5, 20, 11, 2, 253, 251, 3, 2, 2, 2, 253, 252, 3, 2, 2, 2, 254,
	31, 3, 2, 2, 2, 255, 256, 7, 43, 2, 2, 256, 257, 5, 118, 60, 2, 257, 33,
	3, 2, 2, 2, 258, 259, 7, 46, 2, 2, 259, 262, 5, 36, 19, 2, 260, 261, 7,
	10, 2, 2, 261, 263, 5, 36, 19, 2, 262, 260, 3, 2, 2, 2, 262, 263, 3, 2,
	2, 2, 263, 35, 3, 2, 2, 2, 264, 270, 5, 82, 42, 2, 265, 270, 5, 66, 34,
	2, 266, 270, 5, 68, 35, 2, 267, 270, 5, 100, 51, 2, 268, 270, 5, 96, 49,
	2, 269, 264, 3, 2, 2, 2, 269, 265, 3, 2, 2, 2, 269, 266, 3, 2, 2, 2, 269,
	267, 3, 2, 2, 2, 269, 268, 3, 2, 2, 2, 270, 37, 3, 2, 2, 2, 271, 272, 7,
	45, 2, 2, 272, 277, 5, 40, 21, 2, 273, 274, 7, 10, 2, 2, 274, 276, 5, 40,
	21, 2, 275, 273, 3, 2, 2, 2, 276, 279, 3, 2, 2, 2, 277, 275, 3, 2, 2, 2,
	277, 278, 3, 2, 2, 2, 278, 39, 3, 2, 2, 2, 279, 277, 3, 2, 2, 2, 280, 282,
	5, 118, 60, 2, 281, 283, 7, 49, 2, 2, 282, 281, 3, 2, 2, 2, 282, 283, 3,
	2, 2, 2, 283, 41, 3, 2, 2, 2, 284, 285, 7, 48, 2, 2, 285, 303, 5, 54, 28,
	2, 286, 287, 7, 48, 2, 2, 287, 303, 5, 48, 25, 2, 288, 289, 7, 48, 2, 2,
	289, 290, 5, 46, 24, 2, 290, 291, 5, 48, 25, 2, 291, 303, 3, 2, 2, 2, 292,
	293, 7, 48, 2, 2, 293, 294, 5, 46, 24, 2, 294, 295, 5, 52, 27, 2, 295,
	303, 3, 2, 2, 2, 296, 297, 7, 48, 2, 2, 297, 298, 5, 46, 24, 2, 298, 299,
	5, 54, 28, 2, 299, 303, 3, 2, 2, 2, 300, 301, 7, 48, 2, 2, 301, 303, 5,
	46, 24, 2, 302, 284, 3, 2, 2, 2, 302, 286, 3, 2, 2, 2, 302, 288, 3, 2,
	2, 2, 302, 292, 3, 2, 2, 2, 302, 296, 3, 2, 2, 2, 302, 300, 3, 2, 2, 2,
	303, 43, 3, 2, 2, 2, 304, 305, 7, 68, 2, 2, 305, 306, 7, 33, 2, 2, 306,
	307, 5, 118, 60, 2, 307, 45, 3, 2, 2, 2, 308, 313, 5, 44, 23, 2, 309, 310,
	7, 10, 2, 2, 310, 312, 5, 44, 23, 2, 311, 309, 3, 2, 2, 2, 312, 315, 3,
	2, 2, 2, 313, 311, 3, 2, 2, 2, 313, 314, 3, 2, 2, 2, 314, 47, 3, 2, 2,
	2, 315, 313, 3, 2, 2, 2, 316, 317, 7, 60, 2, 2, 317, 322, 5, 50, 26, 2,
	318, 319, 7, 10, 2, 2, 319, 321, 5, 50, 26, 2, 320, 318, 3, 2, 2, 2, 321,
	324, 3, 2, 2, 2, 322, 320, 3, 2, 2, 2, 322, 323, 3, 2, 2, 2, 323, 49, 3,
	2, 2, 2, 324, 322, 3, 2, 2, 2, 325, 326, 7, 68, 2, 2, 326, 327, 7, 33,
	2, 2, 327, 328, 5, 100, 51, 2, 328, 51, 3, 2, 2, 2, 329, 330, 7, 54, 2,
	2, 330, 338, 5, 44, 23, 2, 331, 332, 7, 54, 2, 2, 332, 335, 7, 68, 2, 2,
	333, 334, 7, 55, 2, 2, 334, 336, 7, 68, 2, 2, 335, 333, 3, 2, 2, 2, 335,
	336, 3, 2, 2, 2, 336, 338, 3, 2, 2, 2, 337, 329, 3, 2, 2, 2, 337, 331,
	3, 2, 2, 2, 338, 53, 3, 2, 2, 2, 339, 340, 7, 56, 2, 2, 340, 341, 7, 57,
	2, 2, 341, 342, 7, 54, 2, 2, 342, 343, 7, 68, 2, 2, 343, 55, 3, 2, 2, 2,
	344, 345, 7, 39, 2, 2, 345, 346, 7, 61, 2, 2, 346, 347, 5, 58, 30, 2, 347,
	348, 7, 64, 2, 2, 348, 350, 5, 60, 31, 2, 349, 351, 5, 62, 32, 2, 350,
	349, 3, 2, 2, 2, 350, 351, 3, 2, 2, 2, 351, 353, 3, 2, 2, 2, 352, 354,
	5, 32, 17, 2, 353, 352, 3, 2, 2, 2, 353, 354, 3, 2, 2, 2, 354, 356, 3,
	2, 2, 2, 355, 357, 5, 64, 33, 2, 356, 355, 3, 2, 2, 2, 356, 357, 3, 2,
	2, 2, 357, 57, 3, 2, 2, 2, 358, 364, 5, 78, 40, 2, 359, 364, 5, 68, 35,
	2, 360, 364, 5, 66, 34, 2, 361, 364, 5, 100, 51, 2, 362, 364, 5, 96, 49,
	2, 363, 358, 3, 2, 2, 2, 363, 359, 3, 2, 2, 2, 363, 360, 3, 2, 2, 2, 363,
	361, 3, 2, 2, 2, 363, 362, 3, 2, 2, 2, 364, 59, 3, 2, 2, 2, 365, 369, 5,
	100, 51, 2, 366, 369, 5, 68, 35, 2, 367, 369, 5, 96, 49, 2, 368, 365, 3,
	2, 2, 2, 368, 366, 3, 2, 2, 2, 368, 367, 3, 2, 2, 2, 369, 61, 3, 2, 2,
	2, 370, 371, 7, 40, 2, 2, 371, 372, 5, 74, 38, 2, 372, 63, 3, 2, 2, 2,
	373, 379, 7, 41, 2, 2, 374, 380, 5, 82, 42, 2, 375, 380, 5, 68, 35, 2,
	376, 380, 5, 66, 34, 2, 377, 380, 5, 96, 49, 2, 378, 380, 5, 102, 52, 2,
	379, 374, 3, 2, 2, 2, 379, 375, 3, 2, 2, 2, 379, 376, 3, 2, 2, 2, 379,
	377, 3, 2, 2, 2, 379, 378, 3, 2, 2, 2, 380, 65, 3, 2, 2, 2, 381, 382, 7,
	67, 2, 2, 382, 386, 7, 68, 2, 2, 383, 384, 7, 67, 2, 2, 384, 386, 5, 110,
	56, 2, 385, 381, 3, 2, 2, 2, 385, 383, 3, 2, 2, 2, 386, 67, 3, 2, 2, 2,
	387, 390, 7, 68, 2, 2, 388, 390, 5, 110, 56, 2, 389, 387, 3, 2, 2, 2, 389,
	388, 3, 2, 2, 2, 390, 69, 3, 2, 2, 2, 391, 399, 5, 72, 37, 2, 392, 399,
	5, 74, 38, 2, 393, 399, 5, 76, 39, 2, 394, 399, 5, 78, 40, 2, 395, 399,
	5, 80, 41, 2, 396, 399, 5, 82, 42, 2, 397, 399, 5, 84, 43, 2, 398, 391,
	3, 2, 2, 2, 398, 392, 3, 2, 2, 2, 398, 393, 3, 2, 2, 2, 398, 394, 3, 2,
	2, 2, 398, 395, 3, 2, 2, 2, 398, 396, 3, 2, 2, 2, 398, 397, 3, 2, 2, 2,
	399, 71, 3, 2, 2, 2, 400, 402, 7, 11, 2, 2, 401, 403, 5, 106, 54, 2, 402,
	401, 3, 2, 2, 2, 402, 403, 3, 2, 2, 2, 403, 404, 3, 2, 2, 2, 404, 405,
	7, 12, 2, 2, 405, 73, 3, 2, 2, 2, 406, 418, 7, 15, 2, 2, 407, 412, 5, 86,
	44, 2, 408, 409, 7, 10, 2, 2, 409, 411, 5, 86, 44, 2, 410, 408, 3, 2, 2,
	2, 411, 414, 3, 2, 2, 2, 412, 410, 3, 2, 2, 2, 412, 413, 3, 2, 2, 2, 413,
	416, 3, 2, 2, 2, 414, 412, 3, 2, 2, 2, 415, 417, 7, 10, 2, 2, 416, 415,
	3, 2, 2, 2, 416, 417, 3, 2, 2, 2, 417, 419, 3, 2, 2, 2, 418, 407, 3, 2,
	2, 2, 418, 419, 3, 2, 2, 2, 419, 420, 3, 2, 2, 2, 420, 421, 7, 16, 2, 2,
	421, 75, 3, 2, 2, 2, 422, 423, 7, 52, 2, 2, 423, 77, 3, 2, 2, 2, 424, 425,
	7, 70, 2, 2, 425, 79, 3, 2, 2, 2, 426, 427, 7, 72, 2, 2, 427, 81, 3, 2,
	2, 2, 428, 429, 7, 71, 2, 2, 429, 83, 3, 2, 2, 2, 430, 431, 9, 3, 2, 2,
	431, 85, 3, 2, 2, 2, 432, 433, 5, 90, 46, 2, 433, 434, 7, 7, 2, 2, 434,
	435, 5, 118, 60, 2, 435, 442, 3, 2, 2, 2, 436, 437, 5, 88, 45, 2, 437,
	438, 7, 7, 2, 2, 438, 439, 5, 118, 60, 2, 439, 442, 3, 2, 2, 2, 440, 442,
	5, 68, 35, 2, 441, 432, 3, 2, 2, 2, 441, 436, 3, 2, 2, 2, 441, 440, 3,
	2, 2, 2, 442, 87, 3, 2, 2, 2, 443, 444, 7, 11, 2, 2, 444, 445, 5, 118,
	60, 2, 445, 446, 7, 12, 2, 2, 446, 89, 3, 2, 2, 2, 447, 453, 7, 68, 2,
	2, 448, 453, 5, 78, 40, 2, 449, 453, 5, 66, 34, 2, 450, 453, 5, 110, 56,
	2, 451, 453, 5, 112, 57, 2, 452, 447, 3, 2, 2, 2, 452, 448, 3, 2, 2, 2,
	452, 449, 3, 2, 2, 2, 452, 450, 3, 2, 2, 2, 452, 451, 3, 2, 2, 2, 453,
	91, 3, 2, 2, 2, 454, 455, 5, 94, 48, 2, 455, 456, 7, 68, 2, 2, 456, 93,
	3, 2, 2, 2, 457, 459, 7, 73, 2, 2, 458, 457, 3, 2, 2, 2, 459, 462, 3, 2,
	2, 2, 460, 458, 3, 2, 2, 2, 460, 461, 3, 2, 2, 2, 461, 95, 3, 2, 2, 2,
	462, 460, 3, 2, 2, 2, 463, 465, 5, 98, 50, 2, 464, 466, 5, 108, 55, 2,
	465, 464, 3, 2, 2, 2, 466, 467, 3, 2, 2, 2, 467, 465, 3, 2, 2, 2, 467,
	468, 3, 2, 2, 2, 468, 97, 3, 2, 2, 2, 469, 475, 5, 68, 35, 2, 470, 475,
	5, 66, 34, 2, 471, 475, 5, 72, 37, 2, 472, 475, 5, 74, 38, 2, 473, 475,
	5, 102, 52, 2, 474, 469, 3, 2, 2, 2, 474, 470, 3, 2, 2, 2, 474, 471, 3,
	2, 2, 2, 474, 472, 3, 2, 2, 2, 474, 473, 3, 2, 2, 2, 475, 99, 3, 2, 2,
	2, 476, 478, 5, 102, 52, 2, 477, 479, 5, 144, 73, 2, 478, 477, 3, 2, 2,
	2, 478, 479, 3, 2, 2, 2, 479, 101, 3, 2, 2, 2, 480, 481, 5, 94, 48, 2,
	481, 482, 5, 104, 53, 2, 482, 484, 7, 13, 2, 2, 483, 485, 5, 106, 54, 2,
	484, 483, 3, 2, 2, 2, 484, 485, 3, 2, 2, 2, 485, 486, 3, 2, 2, 2, 486,
	487, 7, 14, 2, 2, 487, 103, 3, 2, 2, 2, 488, 492, 7, 68, 2, 2, 489, 492,
	5, 110, 56, 2, 490, 492, 5, 112, 57, 2, 491, 488, 3, 2, 2, 2, 491, 489,
	3, 2, 2, 2, 491, 490, 3, 2, 2, 2, 492, 105, 3, 2, 2, 2, 493, 498, 5, 118,
	60, 2, 494, 495, 7, 10, 2, 2, 495, 497, 5, 118, 60, 2, 496, 494, 3, 2,
	2, 2, 497, 500, 3, 2, 2, 2, 498, 496, 3, 2, 2, 2, 498, 499, 3, 2, 2, 2,
	499, 502, 3, 2, 2, 2, 500, 498, 3, 2, 2, 2, 501, 503, 7, 10, 2, 2, 502,
	501, 3, 2, 2, 2, 502, 503, 3, 2, 2, 2, 503, 107, 3, 2, 2, 2, 504, 506,
	5, 144, 73, 2, 505, 504, 3, 2, 2, 2, 505, 506, 3, 2, 2, 2, 506, 507, 3,
	2, 2, 2, 507, 508, 7, 9, 2, 2, 508, 516, 5, 90, 46, 2, 509, 510, 5, 144,
	73, 2, 510, 511, 7, 9, 2, 2, 511, 513, 3, 2, 2, 2, 512, 509, 3, 2, 2, 2,
	512, 513, 3, 2, 2, 2, 513, 514, 3, 2, 2, 2, 514, 516, 5, 88, 45, 2, 515,
	505, 3, 2, 2, 2, 515, 512, 3, 2, 2, 2, 516, 109, 3, 2, 2, 2, 517, 518,
	9, 4, 2, 2, 518, 111, 3, 2, 2, 2, 519, 520, 9, 5, 2, 2, 520, 113, 3, 2,
	2, 2, 521, 522, 5, 116, 59, 2, 522, 523, 7, 32, 2, 2, 523, 524, 5, 116,
	59, 2, 524, 115, 3, 2, 2, 2, 525, 529, 5, 82, 42, 2, 526, 529, 5, 68, 35,
	2, 527, 529, 5, 66, 34, 2, 528, 525, 3, 2, 2, 2, 528, 526, 3, 2, 2, 2,
	528, 527, 3, 2, 2, 2, 529, 117, 3, 2, 2, 2, 530, 531, 8, 60, 1, 2, 531,
	532, 5, 132, 67, 2, 532, 533, 5, 118, 60, 7, 533, 536, 3, 2, 2, 2, 534,
	536, 5, 120, 61, 2, 535, 530, 3, 2, 2, 2, 535, 534, 3, 2, 2, 2, 536, 554,
	3, 2, 2, 2, 537, 538, 12, 6, 2, 2, 538, 539, 5, 136, 69, 2, 539, 540, 5,
	118, 60, 7, 540, 553, 3, 2, 2, 2, 541, 542, 12, 5, 2, 2, 542, 543, 5, 138,
	70, 2, 543, 544, 5, 118, 60, 6, 544, 553, 3, 2, 2, 2, 545, 546, 12, 4,
	2, 2, 546, 548, 7, 34, 2, 2, 547, 549, 5, 118, 60, 2, 548, 547, 3, 2, 2,
	2, 548, 549, 3, 2, 2, 2, 549, 550, 3, 2, 2, 2, 550, 551, 7, 7, 2, 2, 551,
	553, 5, 118, 60, 5, 552, 537, 3, 2, 2, 2, 552, 541, 3, 2, 2, 2, 552, 545,
	3, 2, 2, 2, 553, 556, 3, 2, 2, 2, 554, 552, 3, 2, 2, 2, 554, 555, 3, 2,
	2, 2, 555, 119, 3, 2, 2, 2, 556, 554, 3, 2, 2, 2, 557, 558, 8, 61, 1, 2,
	558, 559, 5, 122, 62, 2, 559, 578, 3, 2, 2, 2, 560, 561, 12, 7, 2, 2, 561,
	562, 5, 126, 64, 2, 562, 563, 5, 120, 61, 8, 563, 577, 3, 2, 2, 2, 564,
	565, 12, 6, 2, 2, 565, 566, 5, 124, 63, 2, 566, 567, 5, 120, 61, 7, 567,
	577, 3, 2, 2, 2, 568, 569, 12, 5, 2, 2, 569, 570, 5, 128, 65, 2, 570, 571,
	5, 120, 61, 6, 571, 577, 3, 2, 2, 2, 572, 573, 12, 4, 2, 2, 573, 574, 5,
	130, 66, 2, 574, 575, 5, 120, 61, 5, 575, 577, 3, 2, 2, 2, 576, 560, 3,
	2, 2, 2, 576, 564, 3, 2, 2, 2, 576, 568, 3, 2, 2, 2, 576, 572, 3, 2, 2,
	2, 577, 580, 3, 2, 2, 2, 578, 576, 3, 2, 2, 2, 578, 579, 3, 2, 2, 2, 579,
	121, 3, 2, 2, 2, 580, 578, 3, 2, 2, 2, 581, 582, 8, 62, 1, 2, 582, 599,
	5, 100, 51, 2, 583, 599, 5, 114, 58, 2, 584, 599, 5, 70, 36, 2, 585, 599,
	5, 68, 35, 2, 586, 599, 5, 96, 49, 2, 587, 599, 5, 66, 34, 2, 588, 592,
	7, 13, 2, 2, 589, 593, 5, 20, 11, 2, 590, 593, 5, 56, 29, 2, 591, 593,
	5, 118, 60, 2, 592, 589, 3, 2, 2, 2, 592, 590, 3, 2, 2, 2, 592, 591, 3,
	2, 2, 2, 593, 594, 3, 2, 2, 2, 594, 596, 7, 14, 2, 2, 595, 597, 5, 144,
	73, 2, 596, 595, 3, 2, 2, 2, 596, 597, 3, 2, 2, 2, 597, 599, 3, 2, 2, 2,
	598, 581, 3, 2, 2, 2, 598, 583, 3, 2, 2, 2, 598, 584, 3, 2, 2, 2, 598,
	585, 3, 2, 2, 2, 598, 586, 3, 2, 2, 2, 598, 587, 3, 2, 2, 2, 598, 588,
	3, 2, 2, 2, 599, 614, 3, 2, 2, 2, 600, 601, 12, 12, 2, 2, 601, 602, 5,
	140, 71, 2, 602, 603, 5, 122, 62, 13, 603, 613, 3, 2, 2, 2, 604, 605, 12,
	11, 2, 2, 605, 606, 5, 142, 72, 2, 606, 607, 5, 122, 62, 12, 607, 613,
	3, 2, 2, 2, 608, 609, 12, 10, 2, 2, 609, 610, 5, 134, 68, 2, 610, 611,
	5, 122, 62, 11, 611, 613, 3, 2, 2, 2, 612, 600, 3, 2, 2, 2, 612, 604, 3,
	2, 2, 2, 612, 608, 3, 2, 2, 2, 613, 616, 3, 2, 2, 2, 614, 612, 3, 2, 2,
	2, 614, 615, 3, 2, 2, 2, 615, 123, 3, 2, 2, 2, 616, 614, 3, 2, 2, 2, 617,
	620, 9, 6, 2, 2, 618, 621, 5, 128, 65, 2, 619, 621, 5, 126, 64, 2, 620,
	618, 3, 2, 2, 2, 620, 619, 3, 2, 2, 2, 621, 125, 3, 2, 2, 2, 622, 623,
	9, 7, 2, 2, 623, 127, 3, 2, 2, 2, 624, 626, 7, 63, 2, 2, 625, 624, 3, 2,
	2, 2, 625, 626, 3, 2, 2, 2, 626, 627, 3, 2, 2, 2, 627, 628, 7, 64, 2, 2,
	628, 129, 3, 2, 2, 2, 629, 631, 7, 63, 2, 2, 630, 629, 3, 2, 2, 2, 630,
	631, 3, 2, 2, 2, 631, 632, 3, 2, 2, 2, 632, 633, 7, 62, 2, 2, 633, 131,
	3, 2, 2, 2, 634, 635, 9, 8, 2, 2, 635, 133, 3, 2, 2, 2, 636, 637, 9, 9,
	2, 2, 637, 135, 3, 2, 2, 2, 638, 639, 7, 30, 2, 2, 639, 137, 3, 2, 2, 2,
	640, 641, 7, 31, 2, 2, 641, 139, 3, 2, 2, 2, 642, 643, 9, 10, 2, 2, 643,
	141, 3, 2, 2, 2, 644, 645, 9, 11, 2, 2, 645, 143, 3, 2, 2, 2, 646, 647,
	7, 34, 2, 2, 647, 145, 3, 2, 2, 2, 68, 149, 164, 172, 176, 187, 191, 199,
	206, 214, 221, 226, 235, 241, 245, 249, 253, 262, 269, 277, 282, 302, 313,
	322, 335, 337, 350, 353, 356, 363, 368, 379, 385, 389, 398, 402, 412, 416,
	418, 441, 452, 460, 467, 474, 478, 484, 491, 498, 502, 505, 512, 515, 528,
	535, 548, 552, 554, 576, 578, 592, 596, 598, 612, 614, 620, 625, 630,
}
var literalNames = []string{
	"", "", "", "", "", "':'", "';'", "'.'", "','", "'['", "']'", "'('", "')'",
	"'{'", "'}'", "'>'", "'<'", "'=='", "'>='", "'<='", "'!='", "'*'", "'/'",
	"'%'", "'+'", "'-'", "'--'", "'++'", "", "", "", "'='", "'?'", "'!~'",
	"'=~'", "'FOR'", "'RETURN'", "'WAITFOR'", "'OPTIONS'", "'TIMEOUT'", "'DISTINCT'",
	"'FILTER'", "'CURRENT'", "'SORT'", "'LIMIT'", "'LET'", "'COLLECT'", "",
	"'NONE'", "'NULL'", "", "'USE'", "'INTO'", "'KEEP'", "'WITH'", "'COUNT'",
	"'ALL'", "'ANY'", "'AGGREGATE'", "'EVENT'", "'LIKE'", "", "'IN'", "'DO'",
	"'WHILE'", "'@'",
}
var symbolicNames = []string{
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

var ruleNames = []string{
	"program", "head", "useExpression", "use", "body", "bodyStatement", "bodyExpression",
	"variableDeclaration", "returnExpression", "forExpression", "forExpressionSource",
	"forExpressionClause", "forExpressionStatement", "forExpressionBody", "forExpressionReturn",
	"filterClause", "limitClause", "limitClauseValue", "sortClause", "sortClauseExpression",
	"collectClause", "collectSelector", "collectGrouping", "collectAggregator",
	"collectAggregateSelector", "collectGroupVariable", "collectCounter", "waitForExpression",
	"waitForEventName", "waitForEventSource", "optionsClause", "timeoutClause",
	"param", "variable", "literal", "arrayLiteral", "objectLiteral", "booleanLiteral",
	"stringLiteral", "floatLiteral", "integerLiteral", "noneLiteral", "propertyAssignment",
	"computedPropertyName", "propertyName", "namespaceIdentifier", "namespace",
	"memberExpression", "memberExpressionSource", "functionCallExpression",
	"functionCall", "functionName", "argumentList", "memberExpressionPath",
	"safeReservedWord", "unsafeReservedWord", "rangeOperator", "rangeOperand",
	"expression", "predicate", "expressionAtom", "arrayOperator", "equalityOperator",
	"inOperator", "likeOperator", "unaryOperator", "regexpOperator", "logicalAndOperator",
	"logicalOrOperator", "multiplicativeOperator", "additiveOperator", "errorOperator",
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
	FqlParserOptions           = 38
	FqlParserTimeout           = 39
	FqlParserDistinct          = 40
	FqlParserFilter            = 41
	FqlParserCurrent           = 42
	FqlParserSort              = 43
	FqlParserLimit             = 44
	FqlParserLet               = 45
	FqlParserCollect           = 46
	FqlParserSortDirection     = 47
	FqlParserNone              = 48
	FqlParserNull              = 49
	FqlParserBooleanLiteral    = 50
	FqlParserUse               = 51
	FqlParserInto              = 52
	FqlParserKeep              = 53
	FqlParserWith              = 54
	FqlParserCount             = 55
	FqlParserAll               = 56
	FqlParserAny               = 57
	FqlParserAggregate         = 58
	FqlParserEvent             = 59
	FqlParserLike              = 60
	FqlParserNot               = 61
	FqlParserIn                = 62
	FqlParserDo                = 63
	FqlParserWhile             = 64
	FqlParserParam             = 65
	FqlParserIdentifier        = 66
	FqlParserIgnoreIdentifier  = 67
	FqlParserStringLiteral     = 68
	FqlParserIntegerLiteral    = 69
	FqlParserFloatLiteral      = 70
	FqlParserNamespaceSegment  = 71
	FqlParserUnknownIdentifier = 72
)

// FqlParser rules.
const (
	FqlParserRULE_program                  = 0
	FqlParserRULE_head                     = 1
	FqlParserRULE_useExpression            = 2
	FqlParserRULE_use                      = 3
	FqlParserRULE_body                     = 4
	FqlParserRULE_bodyStatement            = 5
	FqlParserRULE_bodyExpression           = 6
	FqlParserRULE_variableDeclaration      = 7
	FqlParserRULE_returnExpression         = 8
	FqlParserRULE_forExpression            = 9
	FqlParserRULE_forExpressionSource      = 10
	FqlParserRULE_forExpressionClause      = 11
	FqlParserRULE_forExpressionStatement   = 12
	FqlParserRULE_forExpressionBody        = 13
	FqlParserRULE_forExpressionReturn      = 14
	FqlParserRULE_filterClause             = 15
	FqlParserRULE_limitClause              = 16
	FqlParserRULE_limitClauseValue         = 17
	FqlParserRULE_sortClause               = 18
	FqlParserRULE_sortClauseExpression     = 19
	FqlParserRULE_collectClause            = 20
	FqlParserRULE_collectSelector          = 21
	FqlParserRULE_collectGrouping          = 22
	FqlParserRULE_collectAggregator        = 23
	FqlParserRULE_collectAggregateSelector = 24
	FqlParserRULE_collectGroupVariable     = 25
	FqlParserRULE_collectCounter           = 26
	FqlParserRULE_waitForExpression        = 27
	FqlParserRULE_waitForEventName         = 28
	FqlParserRULE_waitForEventSource       = 29
	FqlParserRULE_optionsClause            = 30
	FqlParserRULE_timeoutClause            = 31
	FqlParserRULE_param                    = 32
	FqlParserRULE_variable                 = 33
	FqlParserRULE_literal                  = 34
	FqlParserRULE_arrayLiteral             = 35
	FqlParserRULE_objectLiteral            = 36
	FqlParserRULE_booleanLiteral           = 37
	FqlParserRULE_stringLiteral            = 38
	FqlParserRULE_floatLiteral             = 39
	FqlParserRULE_integerLiteral           = 40
	FqlParserRULE_noneLiteral              = 41
	FqlParserRULE_propertyAssignment       = 42
	FqlParserRULE_computedPropertyName     = 43
	FqlParserRULE_propertyName             = 44
	FqlParserRULE_namespaceIdentifier      = 45
	FqlParserRULE_namespace                = 46
	FqlParserRULE_memberExpression         = 47
	FqlParserRULE_memberExpressionSource   = 48
	FqlParserRULE_functionCallExpression   = 49
	FqlParserRULE_functionCall             = 50
	FqlParserRULE_functionName             = 51
	FqlParserRULE_argumentList             = 52
	FqlParserRULE_memberExpressionPath     = 53
	FqlParserRULE_safeReservedWord         = 54
	FqlParserRULE_unsafeReservedWord       = 55
	FqlParserRULE_rangeOperator            = 56
	FqlParserRULE_rangeOperand             = 57
	FqlParserRULE_expression               = 58
	FqlParserRULE_predicate                = 59
	FqlParserRULE_expressionAtom           = 60
	FqlParserRULE_arrayOperator            = 61
	FqlParserRULE_equalityOperator         = 62
	FqlParserRULE_inOperator               = 63
	FqlParserRULE_likeOperator             = 64
	FqlParserRULE_unaryOperator            = 65
	FqlParserRULE_regexpOperator           = 66
	FqlParserRULE_logicalAndOperator       = 67
	FqlParserRULE_logicalOrOperator        = 68
	FqlParserRULE_multiplicativeOperator   = 69
	FqlParserRULE_additiveOperator         = 70
	FqlParserRULE_errorOperator            = 71
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
	this := p
	_ = this

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
	p.SetState(147)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 0, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			{
				p.SetState(144)
				p.Head()
			}

		}
		p.SetState(149)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 0, p.GetParserRuleContext())
	}
	{
		p.SetState(150)
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
	this := p
	_ = this

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
		p.SetState(152)
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
	this := p
	_ = this

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
		p.SetState(154)
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
	this := p
	_ = this

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
		p.SetState(156)
		p.Match(FqlParserUse)
	}
	{
		p.SetState(157)
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
	this := p
	_ = this

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
	p.SetState(162)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 1, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			{
				p.SetState(159)
				p.BodyStatement()
			}

		}
		p.SetState(164)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 1, p.GetParserRuleContext())
	}
	{
		p.SetState(165)
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

func (s *BodyStatementContext) VariableDeclaration() IVariableDeclarationContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IVariableDeclarationContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IVariableDeclarationContext)
}

func (s *BodyStatementContext) FunctionCallExpression() IFunctionCallExpressionContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IFunctionCallExpressionContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IFunctionCallExpressionContext)
}

func (s *BodyStatementContext) WaitForExpression() IWaitForExpressionContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IWaitForExpressionContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IWaitForExpressionContext)
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
	this := p
	_ = this

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

	p.SetState(170)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 2, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(167)
			p.VariableDeclaration()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(168)
			p.FunctionCallExpression()
		}

	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(169)
			p.WaitForExpression()
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
	this := p
	_ = this

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

	p.SetState(174)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case FqlParserReturn:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(172)
			p.ReturnExpression()
		}

	case FqlParserFor:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(173)
			p.ForExpression()
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}

	return localctx
}

// IVariableDeclarationContext is an interface to support dynamic dispatch.
type IVariableDeclarationContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetId returns the id token.
	GetId() antlr.Token

	// SetId sets the id token.
	SetId(antlr.Token)

	// IsVariableDeclarationContext differentiates from other interfaces.
	IsVariableDeclarationContext()
}

type VariableDeclarationContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
	id     antlr.Token
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

func (s *VariableDeclarationContext) GetId() antlr.Token { return s.id }

func (s *VariableDeclarationContext) SetId(v antlr.Token) { s.id = v }

func (s *VariableDeclarationContext) Let() antlr.TerminalNode {
	return s.GetToken(FqlParserLet, 0)
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

func (s *VariableDeclarationContext) Identifier() antlr.TerminalNode {
	return s.GetToken(FqlParserIdentifier, 0)
}

func (s *VariableDeclarationContext) IgnoreIdentifier() antlr.TerminalNode {
	return s.GetToken(FqlParserIgnoreIdentifier, 0)
}

func (s *VariableDeclarationContext) SafeReservedWord() ISafeReservedWordContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ISafeReservedWordContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ISafeReservedWordContext)
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
	this := p
	_ = this

	localctx = NewVariableDeclarationContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 14, FqlParserRULE_variableDeclaration)
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
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 4, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(176)
			p.Match(FqlParserLet)
		}
		{
			p.SetState(177)

			var _lt = p.GetTokenStream().LT(1)

			localctx.(*VariableDeclarationContext).id = _lt

			_la = p.GetTokenStream().LA(1)

			if !(_la == FqlParserIdentifier || _la == FqlParserIgnoreIdentifier) {
				var _ri = p.GetErrorHandler().RecoverInline(p)

				localctx.(*VariableDeclarationContext).id = _ri
			} else {
				p.GetErrorHandler().ReportMatch(p)
				p.Consume()
			}
		}
		{
			p.SetState(178)
			p.Match(FqlParserAssign)
		}
		{
			p.SetState(179)
			p.expression(0)
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(180)
			p.Match(FqlParserLet)
		}
		{
			p.SetState(181)
			p.SafeReservedWord()
		}
		{
			p.SetState(182)
			p.Match(FqlParserAssign)
		}
		{
			p.SetState(183)
			p.expression(0)
		}

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
	this := p
	_ = this

	localctx = NewReturnExpressionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 16, FqlParserRULE_returnExpression)

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
		p.SetState(187)
		p.Match(FqlParserReturn)
	}
	p.SetState(189)
	p.GetErrorHandler().Sync(p)

	if p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 5, p.GetParserRuleContext()) == 1 {
		{
			p.SetState(188)
			p.Match(FqlParserDistinct)
		}

	}
	{
		p.SetState(191)
		p.expression(0)
	}

	return localctx
}

// IForExpressionContext is an interface to support dynamic dispatch.
type IForExpressionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetValueVariable returns the valueVariable token.
	GetValueVariable() antlr.Token

	// GetCounterVariable returns the counterVariable token.
	GetCounterVariable() antlr.Token

	// SetValueVariable sets the valueVariable token.
	SetValueVariable(antlr.Token)

	// SetCounterVariable sets the counterVariable token.
	SetCounterVariable(antlr.Token)

	// IsForExpressionContext differentiates from other interfaces.
	IsForExpressionContext()
}

type ForExpressionContext struct {
	*antlr.BaseParserRuleContext
	parser          antlr.Parser
	valueVariable   antlr.Token
	counterVariable antlr.Token
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

func (s *ForExpressionContext) GetValueVariable() antlr.Token { return s.valueVariable }

func (s *ForExpressionContext) GetCounterVariable() antlr.Token { return s.counterVariable }

func (s *ForExpressionContext) SetValueVariable(v antlr.Token) { s.valueVariable = v }

func (s *ForExpressionContext) SetCounterVariable(v antlr.Token) { s.counterVariable = v }

func (s *ForExpressionContext) For() antlr.TerminalNode {
	return s.GetToken(FqlParserFor, 0)
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

func (s *ForExpressionContext) AllIdentifier() []antlr.TerminalNode {
	return s.GetTokens(FqlParserIdentifier)
}

func (s *ForExpressionContext) Identifier(i int) antlr.TerminalNode {
	return s.GetToken(FqlParserIdentifier, i)
}

func (s *ForExpressionContext) IgnoreIdentifier() antlr.TerminalNode {
	return s.GetToken(FqlParserIgnoreIdentifier, 0)
}

func (s *ForExpressionContext) Comma() antlr.TerminalNode {
	return s.GetToken(FqlParserComma, 0)
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
	this := p
	_ = this

	localctx = NewForExpressionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 18, FqlParserRULE_forExpression)
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

	p.SetState(224)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 10, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(193)
			p.Match(FqlParserFor)
		}
		{
			p.SetState(194)

			var _lt = p.GetTokenStream().LT(1)

			localctx.(*ForExpressionContext).valueVariable = _lt

			_la = p.GetTokenStream().LA(1)

			if !(_la == FqlParserIdentifier || _la == FqlParserIgnoreIdentifier) {
				var _ri = p.GetErrorHandler().RecoverInline(p)

				localctx.(*ForExpressionContext).valueVariable = _ri
			} else {
				p.GetErrorHandler().ReportMatch(p)
				p.Consume()
			}
		}
		p.SetState(197)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		if _la == FqlParserComma {
			{
				p.SetState(195)
				p.Match(FqlParserComma)
			}
			{
				p.SetState(196)

				var _m = p.Match(FqlParserIdentifier)

				localctx.(*ForExpressionContext).counterVariable = _m
			}

		}
		{
			p.SetState(199)
			p.Match(FqlParserIn)
		}
		{
			p.SetState(200)
			p.ForExpressionSource()
		}
		p.SetState(204)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 7, p.GetParserRuleContext())

		for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
			if _alt == 1 {
				{
					p.SetState(201)
					p.ForExpressionBody()
				}

			}
			p.SetState(206)
			p.GetErrorHandler().Sync(p)
			_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 7, p.GetParserRuleContext())
		}
		{
			p.SetState(207)
			p.ForExpressionReturn()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(209)
			p.Match(FqlParserFor)
		}
		{
			p.SetState(210)

			var _lt = p.GetTokenStream().LT(1)

			localctx.(*ForExpressionContext).counterVariable = _lt

			_la = p.GetTokenStream().LA(1)

			if !(_la == FqlParserIdentifier || _la == FqlParserIgnoreIdentifier) {
				var _ri = p.GetErrorHandler().RecoverInline(p)

				localctx.(*ForExpressionContext).counterVariable = _ri
			} else {
				p.GetErrorHandler().ReportMatch(p)
				p.Consume()
			}
		}
		p.SetState(212)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		if _la == FqlParserDo {
			{
				p.SetState(211)
				p.Match(FqlParserDo)
			}

		}
		{
			p.SetState(214)
			p.Match(FqlParserWhile)
		}
		{
			p.SetState(215)
			p.expression(0)
		}
		p.SetState(219)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 9, p.GetParserRuleContext())

		for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
			if _alt == 1 {
				{
					p.SetState(216)
					p.ForExpressionBody()
				}

			}
			p.SetState(221)
			p.GetErrorHandler().Sync(p)
			_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 9, p.GetParserRuleContext())
		}
		{
			p.SetState(222)
			p.ForExpressionReturn()
		}

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
	this := p
	_ = this

	localctx = NewForExpressionSourceContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 20, FqlParserRULE_forExpressionSource)

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

	p.SetState(233)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 11, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(226)
			p.FunctionCallExpression()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(227)
			p.ArrayLiteral()
		}

	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(228)
			p.ObjectLiteral()
		}

	case 4:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(229)
			p.Variable()
		}

	case 5:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(230)
			p.MemberExpression()
		}

	case 6:
		p.EnterOuterAlt(localctx, 6)
		{
			p.SetState(231)
			p.RangeOperator()
		}

	case 7:
		p.EnterOuterAlt(localctx, 7)
		{
			p.SetState(232)
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
	this := p
	_ = this

	localctx = NewForExpressionClauseContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 22, FqlParserRULE_forExpressionClause)

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

	p.SetState(239)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case FqlParserLimit:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(235)
			p.LimitClause()
		}

	case FqlParserSort:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(236)
			p.SortClause()
		}

	case FqlParserFilter:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(237)
			p.FilterClause()
		}

	case FqlParserCollect:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(238)
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
	this := p
	_ = this

	localctx = NewForExpressionStatementContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 24, FqlParserRULE_forExpressionStatement)

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

	p.SetState(243)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 13, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(241)
			p.VariableDeclaration()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(242)
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
	this := p
	_ = this

	localctx = NewForExpressionBodyContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 26, FqlParserRULE_forExpressionBody)

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

	p.SetState(247)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 14, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(245)
			p.ForExpressionStatement()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(246)
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
	this := p
	_ = this

	localctx = NewForExpressionReturnContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 28, FqlParserRULE_forExpressionReturn)

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

	p.SetState(251)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case FqlParserReturn:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(249)
			p.ReturnExpression()
		}

	case FqlParserFor:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(250)
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
	this := p
	_ = this

	localctx = NewFilterClauseContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 30, FqlParserRULE_filterClause)

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
		p.SetState(253)
		p.Match(FqlParserFilter)
	}
	{
		p.SetState(254)
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
	this := p
	_ = this

	localctx = NewLimitClauseContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 32, FqlParserRULE_limitClause)
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
		p.SetState(256)
		p.Match(FqlParserLimit)
	}
	{
		p.SetState(257)
		p.LimitClauseValue()
	}
	p.SetState(260)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == FqlParserComma {
		{
			p.SetState(258)
			p.Match(FqlParserComma)
		}
		{
			p.SetState(259)
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

func (s *LimitClauseValueContext) IntegerLiteral() IIntegerLiteralContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IIntegerLiteralContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IIntegerLiteralContext)
}

func (s *LimitClauseValueContext) Param() IParamContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IParamContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IParamContext)
}

func (s *LimitClauseValueContext) Variable() IVariableContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IVariableContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IVariableContext)
}

func (s *LimitClauseValueContext) FunctionCallExpression() IFunctionCallExpressionContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IFunctionCallExpressionContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IFunctionCallExpressionContext)
}

func (s *LimitClauseValueContext) MemberExpression() IMemberExpressionContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IMemberExpressionContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IMemberExpressionContext)
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
	this := p
	_ = this

	localctx = NewLimitClauseValueContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 34, FqlParserRULE_limitClauseValue)

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
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 17, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(262)
			p.IntegerLiteral()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(263)
			p.Param()
		}

	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(264)
			p.Variable()
		}

	case 4:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(265)
			p.FunctionCallExpression()
		}

	case 5:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(266)
			p.MemberExpression()
		}

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
	this := p
	_ = this

	localctx = NewSortClauseContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 36, FqlParserRULE_sortClause)
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
		p.SetState(269)
		p.Match(FqlParserSort)
	}
	{
		p.SetState(270)
		p.SortClauseExpression()
	}
	p.SetState(275)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == FqlParserComma {
		{
			p.SetState(271)
			p.Match(FqlParserComma)
		}
		{
			p.SetState(272)
			p.SortClauseExpression()
		}

		p.SetState(277)
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
	this := p
	_ = this

	localctx = NewSortClauseExpressionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 38, FqlParserRULE_sortClauseExpression)

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
		p.SetState(278)
		p.expression(0)
	}
	p.SetState(280)
	p.GetErrorHandler().Sync(p)

	if p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 19, p.GetParserRuleContext()) == 1 {
		{
			p.SetState(279)
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
	this := p
	_ = this

	localctx = NewCollectClauseContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 40, FqlParserRULE_collectClause)

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

	p.SetState(300)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 20, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(282)
			p.Match(FqlParserCollect)
		}
		{
			p.SetState(283)
			p.CollectCounter()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(284)
			p.Match(FqlParserCollect)
		}
		{
			p.SetState(285)
			p.CollectAggregator()
		}

	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(286)
			p.Match(FqlParserCollect)
		}
		{
			p.SetState(287)
			p.CollectGrouping()
		}
		{
			p.SetState(288)
			p.CollectAggregator()
		}

	case 4:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(290)
			p.Match(FqlParserCollect)
		}
		{
			p.SetState(291)
			p.CollectGrouping()
		}
		{
			p.SetState(292)
			p.CollectGroupVariable()
		}

	case 5:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(294)
			p.Match(FqlParserCollect)
		}
		{
			p.SetState(295)
			p.CollectGrouping()
		}
		{
			p.SetState(296)
			p.CollectCounter()
		}

	case 6:
		p.EnterOuterAlt(localctx, 6)
		{
			p.SetState(298)
			p.Match(FqlParserCollect)
		}
		{
			p.SetState(299)
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
	this := p
	_ = this

	localctx = NewCollectSelectorContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 42, FqlParserRULE_collectSelector)

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
		p.SetState(302)
		p.Match(FqlParserIdentifier)
	}
	{
		p.SetState(303)
		p.Match(FqlParserAssign)
	}
	{
		p.SetState(304)
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
	this := p
	_ = this

	localctx = NewCollectGroupingContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 44, FqlParserRULE_collectGrouping)
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
		p.SetState(306)
		p.CollectSelector()
	}
	p.SetState(311)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == FqlParserComma {
		{
			p.SetState(307)
			p.Match(FqlParserComma)
		}
		{
			p.SetState(308)
			p.CollectSelector()
		}

		p.SetState(313)
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
	this := p
	_ = this

	localctx = NewCollectAggregatorContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 46, FqlParserRULE_collectAggregator)
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
		p.SetState(314)
		p.Match(FqlParserAggregate)
	}
	{
		p.SetState(315)
		p.CollectAggregateSelector()
	}
	p.SetState(320)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == FqlParserComma {
		{
			p.SetState(316)
			p.Match(FqlParserComma)
		}
		{
			p.SetState(317)
			p.CollectAggregateSelector()
		}

		p.SetState(322)
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
	this := p
	_ = this

	localctx = NewCollectAggregateSelectorContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 48, FqlParserRULE_collectAggregateSelector)

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
		p.SetState(323)
		p.Match(FqlParserIdentifier)
	}
	{
		p.SetState(324)
		p.Match(FqlParserAssign)
	}
	{
		p.SetState(325)
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
	this := p
	_ = this

	localctx = NewCollectGroupVariableContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 50, FqlParserRULE_collectGroupVariable)

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

	p.SetState(335)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 24, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(327)
			p.Match(FqlParserInto)
		}
		{
			p.SetState(328)
			p.CollectSelector()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(329)
			p.Match(FqlParserInto)
		}
		{
			p.SetState(330)
			p.Match(FqlParserIdentifier)
		}
		p.SetState(333)
		p.GetErrorHandler().Sync(p)

		if p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 23, p.GetParserRuleContext()) == 1 {
			{
				p.SetState(331)
				p.Match(FqlParserKeep)
			}
			{
				p.SetState(332)
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
	this := p
	_ = this

	localctx = NewCollectCounterContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 52, FqlParserRULE_collectCounter)

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
		p.Match(FqlParserWith)
	}
	{
		p.SetState(338)
		p.Match(FqlParserCount)
	}
	{
		p.SetState(339)
		p.Match(FqlParserInto)
	}
	{
		p.SetState(340)
		p.Match(FqlParserIdentifier)
	}

	return localctx
}

// IWaitForExpressionContext is an interface to support dynamic dispatch.
type IWaitForExpressionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsWaitForExpressionContext differentiates from other interfaces.
	IsWaitForExpressionContext()
}

type WaitForExpressionContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyWaitForExpressionContext() *WaitForExpressionContext {
	var p = new(WaitForExpressionContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = FqlParserRULE_waitForExpression
	return p
}

func (*WaitForExpressionContext) IsWaitForExpressionContext() {}

func NewWaitForExpressionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *WaitForExpressionContext {
	var p = new(WaitForExpressionContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_waitForExpression

	return p
}

func (s *WaitForExpressionContext) GetParser() antlr.Parser { return s.parser }

func (s *WaitForExpressionContext) Waitfor() antlr.TerminalNode {
	return s.GetToken(FqlParserWaitfor, 0)
}

func (s *WaitForExpressionContext) Event() antlr.TerminalNode {
	return s.GetToken(FqlParserEvent, 0)
}

func (s *WaitForExpressionContext) WaitForEventName() IWaitForEventNameContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IWaitForEventNameContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IWaitForEventNameContext)
}

func (s *WaitForExpressionContext) In() antlr.TerminalNode {
	return s.GetToken(FqlParserIn, 0)
}

func (s *WaitForExpressionContext) WaitForEventSource() IWaitForEventSourceContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IWaitForEventSourceContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IWaitForEventSourceContext)
}

func (s *WaitForExpressionContext) OptionsClause() IOptionsClauseContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IOptionsClauseContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IOptionsClauseContext)
}

func (s *WaitForExpressionContext) FilterClause() IFilterClauseContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IFilterClauseContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IFilterClauseContext)
}

func (s *WaitForExpressionContext) TimeoutClause() ITimeoutClauseContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ITimeoutClauseContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ITimeoutClauseContext)
}

func (s *WaitForExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *WaitForExpressionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *WaitForExpressionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterWaitForExpression(s)
	}
}

func (s *WaitForExpressionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitWaitForExpression(s)
	}
}

func (s *WaitForExpressionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitWaitForExpression(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) WaitForExpression() (localctx IWaitForExpressionContext) {
	this := p
	_ = this

	localctx = NewWaitForExpressionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 54, FqlParserRULE_waitForExpression)

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
		p.Match(FqlParserWaitfor)
	}
	{
		p.SetState(343)
		p.Match(FqlParserEvent)
	}
	{
		p.SetState(344)
		p.WaitForEventName()
	}
	{
		p.SetState(345)
		p.Match(FqlParserIn)
	}
	{
		p.SetState(346)
		p.WaitForEventSource()
	}
	p.SetState(348)
	p.GetErrorHandler().Sync(p)

	if p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 25, p.GetParserRuleContext()) == 1 {
		{
			p.SetState(347)
			p.OptionsClause()
		}

	}
	p.SetState(351)
	p.GetErrorHandler().Sync(p)

	if p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 26, p.GetParserRuleContext()) == 1 {
		{
			p.SetState(350)
			p.FilterClause()
		}

	}
	p.SetState(354)
	p.GetErrorHandler().Sync(p)

	if p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 27, p.GetParserRuleContext()) == 1 {
		{
			p.SetState(353)
			p.TimeoutClause()
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
	this := p
	_ = this

	localctx = NewWaitForEventNameContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 56, FqlParserRULE_waitForEventName)

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

	p.SetState(361)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 28, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(356)
			p.StringLiteral()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(357)
			p.Variable()
		}

	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(358)
			p.Param()
		}

	case 4:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(359)
			p.FunctionCallExpression()
		}

	case 5:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(360)
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

func (s *WaitForEventSourceContext) FunctionCallExpression() IFunctionCallExpressionContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IFunctionCallExpressionContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IFunctionCallExpressionContext)
}

func (s *WaitForEventSourceContext) Variable() IVariableContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IVariableContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IVariableContext)
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
	this := p
	_ = this

	localctx = NewWaitForEventSourceContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 58, FqlParserRULE_waitForEventSource)

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

	p.SetState(366)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 29, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(363)
			p.FunctionCallExpression()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(364)
			p.Variable()
		}

	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(365)
			p.MemberExpression()
		}

	}

	return localctx
}

// IOptionsClauseContext is an interface to support dynamic dispatch.
type IOptionsClauseContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsOptionsClauseContext differentiates from other interfaces.
	IsOptionsClauseContext()
}

type OptionsClauseContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyOptionsClauseContext() *OptionsClauseContext {
	var p = new(OptionsClauseContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = FqlParserRULE_optionsClause
	return p
}

func (*OptionsClauseContext) IsOptionsClauseContext() {}

func NewOptionsClauseContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *OptionsClauseContext {
	var p = new(OptionsClauseContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_optionsClause

	return p
}

func (s *OptionsClauseContext) GetParser() antlr.Parser { return s.parser }

func (s *OptionsClauseContext) Options() antlr.TerminalNode {
	return s.GetToken(FqlParserOptions, 0)
}

func (s *OptionsClauseContext) ObjectLiteral() IObjectLiteralContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IObjectLiteralContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IObjectLiteralContext)
}

func (s *OptionsClauseContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *OptionsClauseContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *OptionsClauseContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterOptionsClause(s)
	}
}

func (s *OptionsClauseContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitOptionsClause(s)
	}
}

func (s *OptionsClauseContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitOptionsClause(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) OptionsClause() (localctx IOptionsClauseContext) {
	this := p
	_ = this

	localctx = NewOptionsClauseContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 60, FqlParserRULE_optionsClause)

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
		p.Match(FqlParserOptions)
	}
	{
		p.SetState(369)
		p.ObjectLiteral()
	}

	return localctx
}

// ITimeoutClauseContext is an interface to support dynamic dispatch.
type ITimeoutClauseContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsTimeoutClauseContext differentiates from other interfaces.
	IsTimeoutClauseContext()
}

type TimeoutClauseContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyTimeoutClauseContext() *TimeoutClauseContext {
	var p = new(TimeoutClauseContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = FqlParserRULE_timeoutClause
	return p
}

func (*TimeoutClauseContext) IsTimeoutClauseContext() {}

func NewTimeoutClauseContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TimeoutClauseContext {
	var p = new(TimeoutClauseContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_timeoutClause

	return p
}

func (s *TimeoutClauseContext) GetParser() antlr.Parser { return s.parser }

func (s *TimeoutClauseContext) Timeout() antlr.TerminalNode {
	return s.GetToken(FqlParserTimeout, 0)
}

func (s *TimeoutClauseContext) IntegerLiteral() IIntegerLiteralContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IIntegerLiteralContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IIntegerLiteralContext)
}

func (s *TimeoutClauseContext) Variable() IVariableContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IVariableContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IVariableContext)
}

func (s *TimeoutClauseContext) Param() IParamContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IParamContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IParamContext)
}

func (s *TimeoutClauseContext) MemberExpression() IMemberExpressionContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IMemberExpressionContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IMemberExpressionContext)
}

func (s *TimeoutClauseContext) FunctionCall() IFunctionCallContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IFunctionCallContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IFunctionCallContext)
}

func (s *TimeoutClauseContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TimeoutClauseContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *TimeoutClauseContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterTimeoutClause(s)
	}
}

func (s *TimeoutClauseContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitTimeoutClause(s)
	}
}

func (s *TimeoutClauseContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitTimeoutClause(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) TimeoutClause() (localctx ITimeoutClauseContext) {
	this := p
	_ = this

	localctx = NewTimeoutClauseContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 62, FqlParserRULE_timeoutClause)

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
		p.SetState(371)
		p.Match(FqlParserTimeout)
	}
	p.SetState(377)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 30, p.GetParserRuleContext()) {
	case 1:
		{
			p.SetState(372)
			p.IntegerLiteral()
		}

	case 2:
		{
			p.SetState(373)
			p.Variable()
		}

	case 3:
		{
			p.SetState(374)
			p.Param()
		}

	case 4:
		{
			p.SetState(375)
			p.MemberExpression()
		}

	case 5:
		{
			p.SetState(376)
			p.FunctionCall()
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

func (s *ParamContext) SafeReservedWord() ISafeReservedWordContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ISafeReservedWordContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ISafeReservedWordContext)
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
	this := p
	_ = this

	localctx = NewParamContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 64, FqlParserRULE_param)

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

	p.SetState(383)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 31, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(379)
			p.Match(FqlParserParam)
		}
		{
			p.SetState(380)
			p.Match(FqlParserIdentifier)
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(381)
			p.Match(FqlParserParam)
		}
		{
			p.SetState(382)
			p.SafeReservedWord()
		}

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

func (s *VariableContext) SafeReservedWord() ISafeReservedWordContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ISafeReservedWordContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ISafeReservedWordContext)
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
	this := p
	_ = this

	localctx = NewVariableContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 66, FqlParserRULE_variable)

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

	p.SetState(387)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case FqlParserIdentifier:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(385)
			p.Match(FqlParserIdentifier)
		}

	case FqlParserAnd, FqlParserOr, FqlParserOptions, FqlParserTimeout, FqlParserDistinct, FqlParserFilter, FqlParserCurrent, FqlParserSort, FqlParserLimit, FqlParserCollect, FqlParserSortDirection, FqlParserInto, FqlParserKeep, FqlParserWith, FqlParserCount, FqlParserAll, FqlParserAny, FqlParserAggregate, FqlParserEvent:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(386)
			p.SafeReservedWord()
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}

	return localctx
}

// ILiteralContext is an interface to support dynamic dispatch.
type ILiteralContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsLiteralContext differentiates from other interfaces.
	IsLiteralContext()
}

type LiteralContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyLiteralContext() *LiteralContext {
	var p = new(LiteralContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = FqlParserRULE_literal
	return p
}

func (*LiteralContext) IsLiteralContext() {}

func NewLiteralContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *LiteralContext {
	var p = new(LiteralContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_literal

	return p
}

func (s *LiteralContext) GetParser() antlr.Parser { return s.parser }

func (s *LiteralContext) ArrayLiteral() IArrayLiteralContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IArrayLiteralContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IArrayLiteralContext)
}

func (s *LiteralContext) ObjectLiteral() IObjectLiteralContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IObjectLiteralContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IObjectLiteralContext)
}

func (s *LiteralContext) BooleanLiteral() IBooleanLiteralContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IBooleanLiteralContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IBooleanLiteralContext)
}

func (s *LiteralContext) StringLiteral() IStringLiteralContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IStringLiteralContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IStringLiteralContext)
}

func (s *LiteralContext) FloatLiteral() IFloatLiteralContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IFloatLiteralContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IFloatLiteralContext)
}

func (s *LiteralContext) IntegerLiteral() IIntegerLiteralContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IIntegerLiteralContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IIntegerLiteralContext)
}

func (s *LiteralContext) NoneLiteral() INoneLiteralContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*INoneLiteralContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(INoneLiteralContext)
}

func (s *LiteralContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *LiteralContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *LiteralContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterLiteral(s)
	}
}

func (s *LiteralContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitLiteral(s)
	}
}

func (s *LiteralContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitLiteral(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) Literal() (localctx ILiteralContext) {
	this := p
	_ = this

	localctx = NewLiteralContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 68, FqlParserRULE_literal)

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

	p.SetState(396)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case FqlParserOpenBracket:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(389)
			p.ArrayLiteral()
		}

	case FqlParserOpenBrace:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(390)
			p.ObjectLiteral()
		}

	case FqlParserBooleanLiteral:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(391)
			p.BooleanLiteral()
		}

	case FqlParserStringLiteral:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(392)
			p.StringLiteral()
		}

	case FqlParserFloatLiteral:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(393)
			p.FloatLiteral()
		}

	case FqlParserIntegerLiteral:
		p.EnterOuterAlt(localctx, 6)
		{
			p.SetState(394)
			p.IntegerLiteral()
		}

	case FqlParserNone, FqlParserNull:
		p.EnterOuterAlt(localctx, 7)
		{
			p.SetState(395)
			p.NoneLiteral()
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

func (s *ArrayLiteralContext) ArgumentList() IArgumentListContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IArgumentListContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IArgumentListContext)
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
	this := p
	_ = this

	localctx = NewArrayLiteralContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 70, FqlParserRULE_arrayLiteral)
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

	if (((_la-9)&-(0x1f+1)) == 0 && ((1<<uint((_la-9)))&((1<<(FqlParserOpenBracket-9))|(1<<(FqlParserOpenParen-9))|(1<<(FqlParserOpenBrace-9))|(1<<(FqlParserPlus-9))|(1<<(FqlParserMinus-9))|(1<<(FqlParserAnd-9))|(1<<(FqlParserOr-9))|(1<<(FqlParserFor-9))|(1<<(FqlParserReturn-9))|(1<<(FqlParserWaitfor-9))|(1<<(FqlParserOptions-9))|(1<<(FqlParserTimeout-9))|(1<<(FqlParserDistinct-9)))) != 0) || (((_la-41)&-(0x1f+1)) == 0 && ((1<<uint((_la-41)))&((1<<(FqlParserFilter-41))|(1<<(FqlParserCurrent-41))|(1<<(FqlParserSort-41))|(1<<(FqlParserLimit-41))|(1<<(FqlParserLet-41))|(1<<(FqlParserCollect-41))|(1<<(FqlParserSortDirection-41))|(1<<(FqlParserNone-41))|(1<<(FqlParserNull-41))|(1<<(FqlParserBooleanLiteral-41))|(1<<(FqlParserUse-41))|(1<<(FqlParserInto-41))|(1<<(FqlParserKeep-41))|(1<<(FqlParserWith-41))|(1<<(FqlParserCount-41))|(1<<(FqlParserAll-41))|(1<<(FqlParserAny-41))|(1<<(FqlParserAggregate-41))|(1<<(FqlParserEvent-41))|(1<<(FqlParserLike-41))|(1<<(FqlParserNot-41))|(1<<(FqlParserIn-41))|(1<<(FqlParserDo-41))|(1<<(FqlParserWhile-41))|(1<<(FqlParserParam-41))|(1<<(FqlParserIdentifier-41))|(1<<(FqlParserStringLiteral-41))|(1<<(FqlParserIntegerLiteral-41))|(1<<(FqlParserFloatLiteral-41))|(1<<(FqlParserNamespaceSegment-41)))) != 0) {
		{
			p.SetState(399)
			p.ArgumentList()
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
	this := p
	_ = this

	localctx = NewObjectLiteralContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 72, FqlParserRULE_objectLiteral)
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
	p.SetState(416)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if (((_la-9)&-(0x1f+1)) == 0 && ((1<<uint((_la-9)))&((1<<(FqlParserOpenBracket-9))|(1<<(FqlParserAnd-9))|(1<<(FqlParserOr-9))|(1<<(FqlParserFor-9))|(1<<(FqlParserReturn-9))|(1<<(FqlParserWaitfor-9))|(1<<(FqlParserOptions-9))|(1<<(FqlParserTimeout-9))|(1<<(FqlParserDistinct-9)))) != 0) || (((_la-41)&-(0x1f+1)) == 0 && ((1<<uint((_la-41)))&((1<<(FqlParserFilter-41))|(1<<(FqlParserCurrent-41))|(1<<(FqlParserSort-41))|(1<<(FqlParserLimit-41))|(1<<(FqlParserLet-41))|(1<<(FqlParserCollect-41))|(1<<(FqlParserSortDirection-41))|(1<<(FqlParserNone-41))|(1<<(FqlParserNull-41))|(1<<(FqlParserBooleanLiteral-41))|(1<<(FqlParserUse-41))|(1<<(FqlParserInto-41))|(1<<(FqlParserKeep-41))|(1<<(FqlParserWith-41))|(1<<(FqlParserCount-41))|(1<<(FqlParserAll-41))|(1<<(FqlParserAny-41))|(1<<(FqlParserAggregate-41))|(1<<(FqlParserEvent-41))|(1<<(FqlParserLike-41))|(1<<(FqlParserNot-41))|(1<<(FqlParserIn-41))|(1<<(FqlParserDo-41))|(1<<(FqlParserWhile-41))|(1<<(FqlParserParam-41))|(1<<(FqlParserIdentifier-41))|(1<<(FqlParserStringLiteral-41)))) != 0) {
		{
			p.SetState(405)
			p.PropertyAssignment()
		}
		p.SetState(410)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 35, p.GetParserRuleContext())

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
			_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 35, p.GetParserRuleContext())
		}
		p.SetState(414)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		if _la == FqlParserComma {
			{
				p.SetState(413)
				p.Match(FqlParserComma)
			}

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
	this := p
	_ = this

	localctx = NewBooleanLiteralContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 74, FqlParserRULE_booleanLiteral)

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
	this := p
	_ = this

	localctx = NewStringLiteralContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 76, FqlParserRULE_stringLiteral)

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
	this := p
	_ = this

	localctx = NewFloatLiteralContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 78, FqlParserRULE_floatLiteral)

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
		p.Match(FqlParserFloatLiteral)
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
	this := p
	_ = this

	localctx = NewIntegerLiteralContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 80, FqlParserRULE_integerLiteral)

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
		p.Match(FqlParserIntegerLiteral)
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
	this := p
	_ = this

	localctx = NewNoneLiteralContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 82, FqlParserRULE_noneLiteral)
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

func (s *PropertyAssignmentContext) Variable() IVariableContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IVariableContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IVariableContext)
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
	this := p
	_ = this

	localctx = NewPropertyAssignmentContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 84, FqlParserRULE_propertyAssignment)

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

	p.SetState(439)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 38, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(430)
			p.PropertyName()
		}
		{
			p.SetState(431)
			p.Match(FqlParserColon)
		}
		{
			p.SetState(432)
			p.expression(0)
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(434)
			p.ComputedPropertyName()
		}
		{
			p.SetState(435)
			p.Match(FqlParserColon)
		}
		{
			p.SetState(436)
			p.expression(0)
		}

	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(438)
			p.Variable()
		}

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
	this := p
	_ = this

	localctx = NewComputedPropertyNameContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 86, FqlParserRULE_computedPropertyName)

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
		p.SetState(441)
		p.Match(FqlParserOpenBracket)
	}
	{
		p.SetState(442)
		p.expression(0)
	}
	{
		p.SetState(443)
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

func (s *PropertyNameContext) SafeReservedWord() ISafeReservedWordContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ISafeReservedWordContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ISafeReservedWordContext)
}

func (s *PropertyNameContext) UnsafeReservedWord() IUnsafeReservedWordContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IUnsafeReservedWordContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IUnsafeReservedWordContext)
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
	this := p
	_ = this

	localctx = NewPropertyNameContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 88, FqlParserRULE_propertyName)

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

	p.SetState(450)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case FqlParserIdentifier:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(445)
			p.Match(FqlParserIdentifier)
		}

	case FqlParserStringLiteral:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(446)
			p.StringLiteral()
		}

	case FqlParserParam:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(447)
			p.Param()
		}

	case FqlParserAnd, FqlParserOr, FqlParserOptions, FqlParserTimeout, FqlParserDistinct, FqlParserFilter, FqlParserCurrent, FqlParserSort, FqlParserLimit, FqlParserCollect, FqlParserSortDirection, FqlParserInto, FqlParserKeep, FqlParserWith, FqlParserCount, FqlParserAll, FqlParserAny, FqlParserAggregate, FqlParserEvent:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(448)
			p.SafeReservedWord()
		}

	case FqlParserFor, FqlParserReturn, FqlParserWaitfor, FqlParserLet, FqlParserNone, FqlParserNull, FqlParserBooleanLiteral, FqlParserUse, FqlParserLike, FqlParserNot, FqlParserIn, FqlParserDo, FqlParserWhile:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(449)
			p.UnsafeReservedWord()
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
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
	this := p
	_ = this

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
		p.SetState(452)
		p.Namespace()
	}
	{
		p.SetState(453)
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
	this := p
	_ = this

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
	p.SetState(458)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == FqlParserNamespaceSegment {
		{
			p.SetState(455)
			p.Match(FqlParserNamespaceSegment)
		}

		p.SetState(460)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
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

func (s *MemberExpressionContext) MemberExpressionSource() IMemberExpressionSourceContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IMemberExpressionSourceContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IMemberExpressionSourceContext)
}

func (s *MemberExpressionContext) AllMemberExpressionPath() []IMemberExpressionPathContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IMemberExpressionPathContext)(nil)).Elem())
	var tst = make([]IMemberExpressionPathContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IMemberExpressionPathContext)
		}
	}

	return tst
}

func (s *MemberExpressionContext) MemberExpressionPath(i int) IMemberExpressionPathContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IMemberExpressionPathContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IMemberExpressionPathContext)
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
	this := p
	_ = this

	localctx = NewMemberExpressionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 94, FqlParserRULE_memberExpression)

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
		p.SetState(461)
		p.MemberExpressionSource()
	}
	p.SetState(463)
	p.GetErrorHandler().Sync(p)
	_alt = 1
	for ok := true; ok; ok = _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		switch _alt {
		case 1:
			{
				p.SetState(462)
				p.MemberExpressionPath()
			}

		default:
			panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		}

		p.SetState(465)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 41, p.GetParserRuleContext())
	}

	return localctx
}

// IMemberExpressionSourceContext is an interface to support dynamic dispatch.
type IMemberExpressionSourceContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsMemberExpressionSourceContext differentiates from other interfaces.
	IsMemberExpressionSourceContext()
}

type MemberExpressionSourceContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyMemberExpressionSourceContext() *MemberExpressionSourceContext {
	var p = new(MemberExpressionSourceContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = FqlParserRULE_memberExpressionSource
	return p
}

func (*MemberExpressionSourceContext) IsMemberExpressionSourceContext() {}

func NewMemberExpressionSourceContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *MemberExpressionSourceContext {
	var p = new(MemberExpressionSourceContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_memberExpressionSource

	return p
}

func (s *MemberExpressionSourceContext) GetParser() antlr.Parser { return s.parser }

func (s *MemberExpressionSourceContext) Variable() IVariableContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IVariableContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IVariableContext)
}

func (s *MemberExpressionSourceContext) Param() IParamContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IParamContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IParamContext)
}

func (s *MemberExpressionSourceContext) ArrayLiteral() IArrayLiteralContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IArrayLiteralContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IArrayLiteralContext)
}

func (s *MemberExpressionSourceContext) ObjectLiteral() IObjectLiteralContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IObjectLiteralContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IObjectLiteralContext)
}

func (s *MemberExpressionSourceContext) FunctionCall() IFunctionCallContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IFunctionCallContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IFunctionCallContext)
}

func (s *MemberExpressionSourceContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *MemberExpressionSourceContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *MemberExpressionSourceContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterMemberExpressionSource(s)
	}
}

func (s *MemberExpressionSourceContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitMemberExpressionSource(s)
	}
}

func (s *MemberExpressionSourceContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitMemberExpressionSource(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) MemberExpressionSource() (localctx IMemberExpressionSourceContext) {
	this := p
	_ = this

	localctx = NewMemberExpressionSourceContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 96, FqlParserRULE_memberExpressionSource)

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

	p.SetState(472)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 42, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(467)
			p.Variable()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(468)
			p.Param()
		}

	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(469)
			p.ArrayLiteral()
		}

	case 4:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(470)
			p.ObjectLiteral()
		}

	case 5:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(471)
			p.FunctionCall()
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

func (s *FunctionCallExpressionContext) FunctionCall() IFunctionCallContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IFunctionCallContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IFunctionCallContext)
}

func (s *FunctionCallExpressionContext) ErrorOperator() IErrorOperatorContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IErrorOperatorContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IErrorOperatorContext)
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
	this := p
	_ = this

	localctx = NewFunctionCallExpressionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 98, FqlParserRULE_functionCallExpression)

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
		p.SetState(474)
		p.FunctionCall()
	}
	p.SetState(476)
	p.GetErrorHandler().Sync(p)

	if p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 43, p.GetParserRuleContext()) == 1 {
		{
			p.SetState(475)
			p.ErrorOperator()
		}

	}

	return localctx
}

// IFunctionCallContext is an interface to support dynamic dispatch.
type IFunctionCallContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsFunctionCallContext differentiates from other interfaces.
	IsFunctionCallContext()
}

type FunctionCallContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyFunctionCallContext() *FunctionCallContext {
	var p = new(FunctionCallContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = FqlParserRULE_functionCall
	return p
}

func (*FunctionCallContext) IsFunctionCallContext() {}

func NewFunctionCallContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FunctionCallContext {
	var p = new(FunctionCallContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_functionCall

	return p
}

func (s *FunctionCallContext) GetParser() antlr.Parser { return s.parser }

func (s *FunctionCallContext) Namespace() INamespaceContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*INamespaceContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(INamespaceContext)
}

func (s *FunctionCallContext) FunctionName() IFunctionNameContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IFunctionNameContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IFunctionNameContext)
}

func (s *FunctionCallContext) OpenParen() antlr.TerminalNode {
	return s.GetToken(FqlParserOpenParen, 0)
}

func (s *FunctionCallContext) CloseParen() antlr.TerminalNode {
	return s.GetToken(FqlParserCloseParen, 0)
}

func (s *FunctionCallContext) ArgumentList() IArgumentListContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IArgumentListContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IArgumentListContext)
}

func (s *FunctionCallContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FunctionCallContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *FunctionCallContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterFunctionCall(s)
	}
}

func (s *FunctionCallContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitFunctionCall(s)
	}
}

func (s *FunctionCallContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitFunctionCall(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) FunctionCall() (localctx IFunctionCallContext) {
	this := p
	_ = this

	localctx = NewFunctionCallContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 100, FqlParserRULE_functionCall)
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
		p.SetState(478)
		p.Namespace()
	}
	{
		p.SetState(479)
		p.FunctionName()
	}
	{
		p.SetState(480)
		p.Match(FqlParserOpenParen)
	}
	p.SetState(482)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if (((_la-9)&-(0x1f+1)) == 0 && ((1<<uint((_la-9)))&((1<<(FqlParserOpenBracket-9))|(1<<(FqlParserOpenParen-9))|(1<<(FqlParserOpenBrace-9))|(1<<(FqlParserPlus-9))|(1<<(FqlParserMinus-9))|(1<<(FqlParserAnd-9))|(1<<(FqlParserOr-9))|(1<<(FqlParserFor-9))|(1<<(FqlParserReturn-9))|(1<<(FqlParserWaitfor-9))|(1<<(FqlParserOptions-9))|(1<<(FqlParserTimeout-9))|(1<<(FqlParserDistinct-9)))) != 0) || (((_la-41)&-(0x1f+1)) == 0 && ((1<<uint((_la-41)))&((1<<(FqlParserFilter-41))|(1<<(FqlParserCurrent-41))|(1<<(FqlParserSort-41))|(1<<(FqlParserLimit-41))|(1<<(FqlParserLet-41))|(1<<(FqlParserCollect-41))|(1<<(FqlParserSortDirection-41))|(1<<(FqlParserNone-41))|(1<<(FqlParserNull-41))|(1<<(FqlParserBooleanLiteral-41))|(1<<(FqlParserUse-41))|(1<<(FqlParserInto-41))|(1<<(FqlParserKeep-41))|(1<<(FqlParserWith-41))|(1<<(FqlParserCount-41))|(1<<(FqlParserAll-41))|(1<<(FqlParserAny-41))|(1<<(FqlParserAggregate-41))|(1<<(FqlParserEvent-41))|(1<<(FqlParserLike-41))|(1<<(FqlParserNot-41))|(1<<(FqlParserIn-41))|(1<<(FqlParserDo-41))|(1<<(FqlParserWhile-41))|(1<<(FqlParserParam-41))|(1<<(FqlParserIdentifier-41))|(1<<(FqlParserStringLiteral-41))|(1<<(FqlParserIntegerLiteral-41))|(1<<(FqlParserFloatLiteral-41))|(1<<(FqlParserNamespaceSegment-41)))) != 0) {
		{
			p.SetState(481)
			p.ArgumentList()
		}

	}
	{
		p.SetState(484)
		p.Match(FqlParserCloseParen)
	}

	return localctx
}

// IFunctionNameContext is an interface to support dynamic dispatch.
type IFunctionNameContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsFunctionNameContext differentiates from other interfaces.
	IsFunctionNameContext()
}

type FunctionNameContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyFunctionNameContext() *FunctionNameContext {
	var p = new(FunctionNameContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = FqlParserRULE_functionName
	return p
}

func (*FunctionNameContext) IsFunctionNameContext() {}

func NewFunctionNameContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FunctionNameContext {
	var p = new(FunctionNameContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_functionName

	return p
}

func (s *FunctionNameContext) GetParser() antlr.Parser { return s.parser }

func (s *FunctionNameContext) Identifier() antlr.TerminalNode {
	return s.GetToken(FqlParserIdentifier, 0)
}

func (s *FunctionNameContext) SafeReservedWord() ISafeReservedWordContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ISafeReservedWordContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ISafeReservedWordContext)
}

func (s *FunctionNameContext) UnsafeReservedWord() IUnsafeReservedWordContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IUnsafeReservedWordContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IUnsafeReservedWordContext)
}

func (s *FunctionNameContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FunctionNameContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *FunctionNameContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterFunctionName(s)
	}
}

func (s *FunctionNameContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitFunctionName(s)
	}
}

func (s *FunctionNameContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitFunctionName(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) FunctionName() (localctx IFunctionNameContext) {
	this := p
	_ = this

	localctx = NewFunctionNameContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 102, FqlParserRULE_functionName)

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

	p.SetState(489)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case FqlParserIdentifier:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(486)
			p.Match(FqlParserIdentifier)
		}

	case FqlParserAnd, FqlParserOr, FqlParserOptions, FqlParserTimeout, FqlParserDistinct, FqlParserFilter, FqlParserCurrent, FqlParserSort, FqlParserLimit, FqlParserCollect, FqlParserSortDirection, FqlParserInto, FqlParserKeep, FqlParserWith, FqlParserCount, FqlParserAll, FqlParserAny, FqlParserAggregate, FqlParserEvent:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(487)
			p.SafeReservedWord()
		}

	case FqlParserFor, FqlParserReturn, FqlParserWaitfor, FqlParserLet, FqlParserNone, FqlParserNull, FqlParserBooleanLiteral, FqlParserUse, FqlParserLike, FqlParserNot, FqlParserIn, FqlParserDo, FqlParserWhile:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(488)
			p.UnsafeReservedWord()
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}

	return localctx
}

// IArgumentListContext is an interface to support dynamic dispatch.
type IArgumentListContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsArgumentListContext differentiates from other interfaces.
	IsArgumentListContext()
}

type ArgumentListContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyArgumentListContext() *ArgumentListContext {
	var p = new(ArgumentListContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = FqlParserRULE_argumentList
	return p
}

func (*ArgumentListContext) IsArgumentListContext() {}

func NewArgumentListContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ArgumentListContext {
	var p = new(ArgumentListContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_argumentList

	return p
}

func (s *ArgumentListContext) GetParser() antlr.Parser { return s.parser }

func (s *ArgumentListContext) AllExpression() []IExpressionContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IExpressionContext)(nil)).Elem())
	var tst = make([]IExpressionContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IExpressionContext)
		}
	}

	return tst
}

func (s *ArgumentListContext) Expression(i int) IExpressionContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExpressionContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *ArgumentListContext) AllComma() []antlr.TerminalNode {
	return s.GetTokens(FqlParserComma)
}

func (s *ArgumentListContext) Comma(i int) antlr.TerminalNode {
	return s.GetToken(FqlParserComma, i)
}

func (s *ArgumentListContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ArgumentListContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ArgumentListContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterArgumentList(s)
	}
}

func (s *ArgumentListContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitArgumentList(s)
	}
}

func (s *ArgumentListContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitArgumentList(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) ArgumentList() (localctx IArgumentListContext) {
	this := p
	_ = this

	localctx = NewArgumentListContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 104, FqlParserRULE_argumentList)
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
		p.SetState(491)
		p.expression(0)
	}
	p.SetState(496)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 46, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			{
				p.SetState(492)
				p.Match(FqlParserComma)
			}
			{
				p.SetState(493)
				p.expression(0)
			}

		}
		p.SetState(498)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 46, p.GetParserRuleContext())
	}
	p.SetState(500)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == FqlParserComma {
		{
			p.SetState(499)
			p.Match(FqlParserComma)
		}

	}

	return localctx
}

// IMemberExpressionPathContext is an interface to support dynamic dispatch.
type IMemberExpressionPathContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsMemberExpressionPathContext differentiates from other interfaces.
	IsMemberExpressionPathContext()
}

type MemberExpressionPathContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyMemberExpressionPathContext() *MemberExpressionPathContext {
	var p = new(MemberExpressionPathContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = FqlParserRULE_memberExpressionPath
	return p
}

func (*MemberExpressionPathContext) IsMemberExpressionPathContext() {}

func NewMemberExpressionPathContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *MemberExpressionPathContext {
	var p = new(MemberExpressionPathContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_memberExpressionPath

	return p
}

func (s *MemberExpressionPathContext) GetParser() antlr.Parser { return s.parser }

func (s *MemberExpressionPathContext) Dot() antlr.TerminalNode {
	return s.GetToken(FqlParserDot, 0)
}

func (s *MemberExpressionPathContext) PropertyName() IPropertyNameContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IPropertyNameContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IPropertyNameContext)
}

func (s *MemberExpressionPathContext) ErrorOperator() IErrorOperatorContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IErrorOperatorContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IErrorOperatorContext)
}

func (s *MemberExpressionPathContext) ComputedPropertyName() IComputedPropertyNameContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IComputedPropertyNameContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IComputedPropertyNameContext)
}

func (s *MemberExpressionPathContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *MemberExpressionPathContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *MemberExpressionPathContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterMemberExpressionPath(s)
	}
}

func (s *MemberExpressionPathContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitMemberExpressionPath(s)
	}
}

func (s *MemberExpressionPathContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitMemberExpressionPath(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) MemberExpressionPath() (localctx IMemberExpressionPathContext) {
	this := p
	_ = this

	localctx = NewMemberExpressionPathContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 106, FqlParserRULE_memberExpressionPath)
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

	p.SetState(513)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 50, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		p.SetState(503)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		if _la == FqlParserQuestionMark {
			{
				p.SetState(502)
				p.ErrorOperator()
			}

		}
		{
			p.SetState(505)
			p.Match(FqlParserDot)
		}
		{
			p.SetState(506)
			p.PropertyName()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		p.SetState(510)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		if _la == FqlParserQuestionMark {
			{
				p.SetState(507)
				p.ErrorOperator()
			}
			{
				p.SetState(508)
				p.Match(FqlParserDot)
			}

		}
		{
			p.SetState(512)
			p.ComputedPropertyName()
		}

	}

	return localctx
}

// ISafeReservedWordContext is an interface to support dynamic dispatch.
type ISafeReservedWordContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsSafeReservedWordContext differentiates from other interfaces.
	IsSafeReservedWordContext()
}

type SafeReservedWordContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptySafeReservedWordContext() *SafeReservedWordContext {
	var p = new(SafeReservedWordContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = FqlParserRULE_safeReservedWord
	return p
}

func (*SafeReservedWordContext) IsSafeReservedWordContext() {}

func NewSafeReservedWordContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SafeReservedWordContext {
	var p = new(SafeReservedWordContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_safeReservedWord

	return p
}

func (s *SafeReservedWordContext) GetParser() antlr.Parser { return s.parser }

func (s *SafeReservedWordContext) And() antlr.TerminalNode {
	return s.GetToken(FqlParserAnd, 0)
}

func (s *SafeReservedWordContext) Or() antlr.TerminalNode {
	return s.GetToken(FqlParserOr, 0)
}

func (s *SafeReservedWordContext) Distinct() antlr.TerminalNode {
	return s.GetToken(FqlParserDistinct, 0)
}

func (s *SafeReservedWordContext) Filter() antlr.TerminalNode {
	return s.GetToken(FqlParserFilter, 0)
}

func (s *SafeReservedWordContext) Sort() antlr.TerminalNode {
	return s.GetToken(FqlParserSort, 0)
}

func (s *SafeReservedWordContext) Limit() antlr.TerminalNode {
	return s.GetToken(FqlParserLimit, 0)
}

func (s *SafeReservedWordContext) Collect() antlr.TerminalNode {
	return s.GetToken(FqlParserCollect, 0)
}

func (s *SafeReservedWordContext) SortDirection() antlr.TerminalNode {
	return s.GetToken(FqlParserSortDirection, 0)
}

func (s *SafeReservedWordContext) Into() antlr.TerminalNode {
	return s.GetToken(FqlParserInto, 0)
}

func (s *SafeReservedWordContext) Keep() antlr.TerminalNode {
	return s.GetToken(FqlParserKeep, 0)
}

func (s *SafeReservedWordContext) With() antlr.TerminalNode {
	return s.GetToken(FqlParserWith, 0)
}

func (s *SafeReservedWordContext) Count() antlr.TerminalNode {
	return s.GetToken(FqlParserCount, 0)
}

func (s *SafeReservedWordContext) All() antlr.TerminalNode {
	return s.GetToken(FqlParserAll, 0)
}

func (s *SafeReservedWordContext) Any() antlr.TerminalNode {
	return s.GetToken(FqlParserAny, 0)
}

func (s *SafeReservedWordContext) Aggregate() antlr.TerminalNode {
	return s.GetToken(FqlParserAggregate, 0)
}

func (s *SafeReservedWordContext) Event() antlr.TerminalNode {
	return s.GetToken(FqlParserEvent, 0)
}

func (s *SafeReservedWordContext) Timeout() antlr.TerminalNode {
	return s.GetToken(FqlParserTimeout, 0)
}

func (s *SafeReservedWordContext) Options() antlr.TerminalNode {
	return s.GetToken(FqlParserOptions, 0)
}

func (s *SafeReservedWordContext) Current() antlr.TerminalNode {
	return s.GetToken(FqlParserCurrent, 0)
}

func (s *SafeReservedWordContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SafeReservedWordContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *SafeReservedWordContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterSafeReservedWord(s)
	}
}

func (s *SafeReservedWordContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitSafeReservedWord(s)
	}
}

func (s *SafeReservedWordContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitSafeReservedWord(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) SafeReservedWord() (localctx ISafeReservedWordContext) {
	this := p
	_ = this

	localctx = NewSafeReservedWordContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 108, FqlParserRULE_safeReservedWord)
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
		p.SetState(515)
		_la = p.GetTokenStream().LA(1)

		if !(((_la-28)&-(0x1f+1)) == 0 && ((1<<uint((_la-28)))&((1<<(FqlParserAnd-28))|(1<<(FqlParserOr-28))|(1<<(FqlParserOptions-28))|(1<<(FqlParserTimeout-28))|(1<<(FqlParserDistinct-28))|(1<<(FqlParserFilter-28))|(1<<(FqlParserCurrent-28))|(1<<(FqlParserSort-28))|(1<<(FqlParserLimit-28))|(1<<(FqlParserCollect-28))|(1<<(FqlParserSortDirection-28))|(1<<(FqlParserInto-28))|(1<<(FqlParserKeep-28))|(1<<(FqlParserWith-28))|(1<<(FqlParserCount-28))|(1<<(FqlParserAll-28))|(1<<(FqlParserAny-28))|(1<<(FqlParserAggregate-28))|(1<<(FqlParserEvent-28)))) != 0) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}

	return localctx
}

// IUnsafeReservedWordContext is an interface to support dynamic dispatch.
type IUnsafeReservedWordContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsUnsafeReservedWordContext differentiates from other interfaces.
	IsUnsafeReservedWordContext()
}

type UnsafeReservedWordContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyUnsafeReservedWordContext() *UnsafeReservedWordContext {
	var p = new(UnsafeReservedWordContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = FqlParserRULE_unsafeReservedWord
	return p
}

func (*UnsafeReservedWordContext) IsUnsafeReservedWordContext() {}

func NewUnsafeReservedWordContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *UnsafeReservedWordContext {
	var p = new(UnsafeReservedWordContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_unsafeReservedWord

	return p
}

func (s *UnsafeReservedWordContext) GetParser() antlr.Parser { return s.parser }

func (s *UnsafeReservedWordContext) Return() antlr.TerminalNode {
	return s.GetToken(FqlParserReturn, 0)
}

func (s *UnsafeReservedWordContext) None() antlr.TerminalNode {
	return s.GetToken(FqlParserNone, 0)
}

func (s *UnsafeReservedWordContext) Null() antlr.TerminalNode {
	return s.GetToken(FqlParserNull, 0)
}

func (s *UnsafeReservedWordContext) Let() antlr.TerminalNode {
	return s.GetToken(FqlParserLet, 0)
}

func (s *UnsafeReservedWordContext) Use() antlr.TerminalNode {
	return s.GetToken(FqlParserUse, 0)
}

func (s *UnsafeReservedWordContext) Waitfor() antlr.TerminalNode {
	return s.GetToken(FqlParserWaitfor, 0)
}

func (s *UnsafeReservedWordContext) While() antlr.TerminalNode {
	return s.GetToken(FqlParserWhile, 0)
}

func (s *UnsafeReservedWordContext) Do() antlr.TerminalNode {
	return s.GetToken(FqlParserDo, 0)
}

func (s *UnsafeReservedWordContext) In() antlr.TerminalNode {
	return s.GetToken(FqlParserIn, 0)
}

func (s *UnsafeReservedWordContext) Like() antlr.TerminalNode {
	return s.GetToken(FqlParserLike, 0)
}

func (s *UnsafeReservedWordContext) Not() antlr.TerminalNode {
	return s.GetToken(FqlParserNot, 0)
}

func (s *UnsafeReservedWordContext) For() antlr.TerminalNode {
	return s.GetToken(FqlParserFor, 0)
}

func (s *UnsafeReservedWordContext) BooleanLiteral() antlr.TerminalNode {
	return s.GetToken(FqlParserBooleanLiteral, 0)
}

func (s *UnsafeReservedWordContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *UnsafeReservedWordContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *UnsafeReservedWordContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterUnsafeReservedWord(s)
	}
}

func (s *UnsafeReservedWordContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitUnsafeReservedWord(s)
	}
}

func (s *UnsafeReservedWordContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitUnsafeReservedWord(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) UnsafeReservedWord() (localctx IUnsafeReservedWordContext) {
	this := p
	_ = this

	localctx = NewUnsafeReservedWordContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 110, FqlParserRULE_unsafeReservedWord)
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
		p.SetState(517)
		_la = p.GetTokenStream().LA(1)

		if !(((_la-35)&-(0x1f+1)) == 0 && ((1<<uint((_la-35)))&((1<<(FqlParserFor-35))|(1<<(FqlParserReturn-35))|(1<<(FqlParserWaitfor-35))|(1<<(FqlParserLet-35))|(1<<(FqlParserNone-35))|(1<<(FqlParserNull-35))|(1<<(FqlParserBooleanLiteral-35))|(1<<(FqlParserUse-35))|(1<<(FqlParserLike-35))|(1<<(FqlParserNot-35))|(1<<(FqlParserIn-35))|(1<<(FqlParserDo-35))|(1<<(FqlParserWhile-35)))) != 0) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}

	return localctx
}

// IRangeOperatorContext is an interface to support dynamic dispatch.
type IRangeOperatorContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetLeft returns the left rule contexts.
	GetLeft() IRangeOperandContext

	// GetRight returns the right rule contexts.
	GetRight() IRangeOperandContext

	// SetLeft sets the left rule contexts.
	SetLeft(IRangeOperandContext)

	// SetRight sets the right rule contexts.
	SetRight(IRangeOperandContext)

	// IsRangeOperatorContext differentiates from other interfaces.
	IsRangeOperatorContext()
}

type RangeOperatorContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
	left   IRangeOperandContext
	right  IRangeOperandContext
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

func (s *RangeOperatorContext) GetLeft() IRangeOperandContext { return s.left }

func (s *RangeOperatorContext) GetRight() IRangeOperandContext { return s.right }

func (s *RangeOperatorContext) SetLeft(v IRangeOperandContext) { s.left = v }

func (s *RangeOperatorContext) SetRight(v IRangeOperandContext) { s.right = v }

func (s *RangeOperatorContext) Range() antlr.TerminalNode {
	return s.GetToken(FqlParserRange, 0)
}

func (s *RangeOperatorContext) AllRangeOperand() []IRangeOperandContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IRangeOperandContext)(nil)).Elem())
	var tst = make([]IRangeOperandContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IRangeOperandContext)
		}
	}

	return tst
}

func (s *RangeOperatorContext) RangeOperand(i int) IRangeOperandContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IRangeOperandContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IRangeOperandContext)
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
	this := p
	_ = this

	localctx = NewRangeOperatorContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 112, FqlParserRULE_rangeOperator)

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
		p.SetState(519)

		var _x = p.RangeOperand()

		localctx.(*RangeOperatorContext).left = _x
	}
	{
		p.SetState(520)
		p.Match(FqlParserRange)
	}
	{
		p.SetState(521)

		var _x = p.RangeOperand()

		localctx.(*RangeOperatorContext).right = _x
	}

	return localctx
}

// IRangeOperandContext is an interface to support dynamic dispatch.
type IRangeOperandContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsRangeOperandContext differentiates from other interfaces.
	IsRangeOperandContext()
}

type RangeOperandContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyRangeOperandContext() *RangeOperandContext {
	var p = new(RangeOperandContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = FqlParserRULE_rangeOperand
	return p
}

func (*RangeOperandContext) IsRangeOperandContext() {}

func NewRangeOperandContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *RangeOperandContext {
	var p = new(RangeOperandContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_rangeOperand

	return p
}

func (s *RangeOperandContext) GetParser() antlr.Parser { return s.parser }

func (s *RangeOperandContext) IntegerLiteral() IIntegerLiteralContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IIntegerLiteralContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IIntegerLiteralContext)
}

func (s *RangeOperandContext) Variable() IVariableContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IVariableContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IVariableContext)
}

func (s *RangeOperandContext) Param() IParamContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IParamContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IParamContext)
}

func (s *RangeOperandContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *RangeOperandContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *RangeOperandContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterRangeOperand(s)
	}
}

func (s *RangeOperandContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitRangeOperand(s)
	}
}

func (s *RangeOperandContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitRangeOperand(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) RangeOperand() (localctx IRangeOperandContext) {
	this := p
	_ = this

	localctx = NewRangeOperandContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 114, FqlParserRULE_rangeOperand)

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

	p.SetState(526)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case FqlParserIntegerLiteral:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(523)
			p.IntegerLiteral()
		}

	case FqlParserAnd, FqlParserOr, FqlParserOptions, FqlParserTimeout, FqlParserDistinct, FqlParserFilter, FqlParserCurrent, FqlParserSort, FqlParserLimit, FqlParserCollect, FqlParserSortDirection, FqlParserInto, FqlParserKeep, FqlParserWith, FqlParserCount, FqlParserAll, FqlParserAny, FqlParserAggregate, FqlParserEvent, FqlParserIdentifier:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(524)
			p.Variable()
		}

	case FqlParserParam:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(525)
			p.Param()
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}

	return localctx
}

// IExpressionContext is an interface to support dynamic dispatch.
type IExpressionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetTernaryOperator returns the ternaryOperator token.
	GetTernaryOperator() antlr.Token

	// SetTernaryOperator sets the ternaryOperator token.
	SetTernaryOperator(antlr.Token)

	// GetLeft returns the left rule contexts.
	GetLeft() IExpressionContext

	// GetCondition returns the condition rule contexts.
	GetCondition() IExpressionContext

	// GetRight returns the right rule contexts.
	GetRight() IExpressionContext

	// GetOnTrue returns the onTrue rule contexts.
	GetOnTrue() IExpressionContext

	// GetOnFalse returns the onFalse rule contexts.
	GetOnFalse() IExpressionContext

	// SetLeft sets the left rule contexts.
	SetLeft(IExpressionContext)

	// SetCondition sets the condition rule contexts.
	SetCondition(IExpressionContext)

	// SetRight sets the right rule contexts.
	SetRight(IExpressionContext)

	// SetOnTrue sets the onTrue rule contexts.
	SetOnTrue(IExpressionContext)

	// SetOnFalse sets the onFalse rule contexts.
	SetOnFalse(IExpressionContext)

	// IsExpressionContext differentiates from other interfaces.
	IsExpressionContext()
}

type ExpressionContext struct {
	*antlr.BaseParserRuleContext
	parser          antlr.Parser
	left            IExpressionContext
	condition       IExpressionContext
	right           IExpressionContext
	ternaryOperator antlr.Token
	onTrue          IExpressionContext
	onFalse         IExpressionContext
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

func (s *ExpressionContext) GetTernaryOperator() antlr.Token { return s.ternaryOperator }

func (s *ExpressionContext) SetTernaryOperator(v antlr.Token) { s.ternaryOperator = v }

func (s *ExpressionContext) GetLeft() IExpressionContext { return s.left }

func (s *ExpressionContext) GetCondition() IExpressionContext { return s.condition }

func (s *ExpressionContext) GetRight() IExpressionContext { return s.right }

func (s *ExpressionContext) GetOnTrue() IExpressionContext { return s.onTrue }

func (s *ExpressionContext) GetOnFalse() IExpressionContext { return s.onFalse }

func (s *ExpressionContext) SetLeft(v IExpressionContext) { s.left = v }

func (s *ExpressionContext) SetCondition(v IExpressionContext) { s.condition = v }

func (s *ExpressionContext) SetRight(v IExpressionContext) { s.right = v }

func (s *ExpressionContext) SetOnTrue(v IExpressionContext) { s.onTrue = v }

func (s *ExpressionContext) SetOnFalse(v IExpressionContext) { s.onFalse = v }

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

func (s *ExpressionContext) Predicate() IPredicateContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IPredicateContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IPredicateContext)
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

func (s *ExpressionContext) Colon() antlr.TerminalNode {
	return s.GetToken(FqlParserColon, 0)
}

func (s *ExpressionContext) QuestionMark() antlr.TerminalNode {
	return s.GetToken(FqlParserQuestionMark, 0)
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
	this := p
	_ = this

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
	p.SetState(533)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 52, p.GetParserRuleContext()) {
	case 1:
		{
			p.SetState(529)
			p.UnaryOperator()
		}
		{
			p.SetState(530)

			var _x = p.expression(5)

			localctx.(*ExpressionContext).right = _x
		}

	case 2:
		{
			p.SetState(532)
			p.predicate(0)
		}

	}
	p.GetParserRuleContext().SetStop(p.GetTokenStream().LT(-1))
	p.SetState(552)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 55, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			if p.GetParseListeners() != nil {
				p.TriggerExitRuleEvent()
			}
			_prevctx = localctx
			p.SetState(550)
			p.GetErrorHandler().Sync(p)
			switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 54, p.GetParserRuleContext()) {
			case 1:
				localctx = NewExpressionContext(p, _parentctx, _parentState)
				localctx.(*ExpressionContext).left = _prevctx
				p.PushNewRecursionContext(localctx, _startState, FqlParserRULE_expression)
				p.SetState(535)

				if !(p.Precpred(p.GetParserRuleContext(), 4)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 4)", ""))
				}
				{
					p.SetState(536)
					p.LogicalAndOperator()
				}
				{
					p.SetState(537)

					var _x = p.expression(5)

					localctx.(*ExpressionContext).right = _x
				}

			case 2:
				localctx = NewExpressionContext(p, _parentctx, _parentState)
				localctx.(*ExpressionContext).left = _prevctx
				p.PushNewRecursionContext(localctx, _startState, FqlParserRULE_expression)
				p.SetState(539)

				if !(p.Precpred(p.GetParserRuleContext(), 3)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 3)", ""))
				}
				{
					p.SetState(540)
					p.LogicalOrOperator()
				}
				{
					p.SetState(541)

					var _x = p.expression(4)

					localctx.(*ExpressionContext).right = _x
				}

			case 3:
				localctx = NewExpressionContext(p, _parentctx, _parentState)
				localctx.(*ExpressionContext).condition = _prevctx
				p.PushNewRecursionContext(localctx, _startState, FqlParserRULE_expression)
				p.SetState(543)

				if !(p.Precpred(p.GetParserRuleContext(), 2)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 2)", ""))
				}
				{
					p.SetState(544)

					var _m = p.Match(FqlParserQuestionMark)

					localctx.(*ExpressionContext).ternaryOperator = _m
				}
				p.SetState(546)
				p.GetErrorHandler().Sync(p)
				_la = p.GetTokenStream().LA(1)

				if (((_la-9)&-(0x1f+1)) == 0 && ((1<<uint((_la-9)))&((1<<(FqlParserOpenBracket-9))|(1<<(FqlParserOpenParen-9))|(1<<(FqlParserOpenBrace-9))|(1<<(FqlParserPlus-9))|(1<<(FqlParserMinus-9))|(1<<(FqlParserAnd-9))|(1<<(FqlParserOr-9))|(1<<(FqlParserFor-9))|(1<<(FqlParserReturn-9))|(1<<(FqlParserWaitfor-9))|(1<<(FqlParserOptions-9))|(1<<(FqlParserTimeout-9))|(1<<(FqlParserDistinct-9)))) != 0) || (((_la-41)&-(0x1f+1)) == 0 && ((1<<uint((_la-41)))&((1<<(FqlParserFilter-41))|(1<<(FqlParserCurrent-41))|(1<<(FqlParserSort-41))|(1<<(FqlParserLimit-41))|(1<<(FqlParserLet-41))|(1<<(FqlParserCollect-41))|(1<<(FqlParserSortDirection-41))|(1<<(FqlParserNone-41))|(1<<(FqlParserNull-41))|(1<<(FqlParserBooleanLiteral-41))|(1<<(FqlParserUse-41))|(1<<(FqlParserInto-41))|(1<<(FqlParserKeep-41))|(1<<(FqlParserWith-41))|(1<<(FqlParserCount-41))|(1<<(FqlParserAll-41))|(1<<(FqlParserAny-41))|(1<<(FqlParserAggregate-41))|(1<<(FqlParserEvent-41))|(1<<(FqlParserLike-41))|(1<<(FqlParserNot-41))|(1<<(FqlParserIn-41))|(1<<(FqlParserDo-41))|(1<<(FqlParserWhile-41))|(1<<(FqlParserParam-41))|(1<<(FqlParserIdentifier-41))|(1<<(FqlParserStringLiteral-41))|(1<<(FqlParserIntegerLiteral-41))|(1<<(FqlParserFloatLiteral-41))|(1<<(FqlParserNamespaceSegment-41)))) != 0) {
					{
						p.SetState(545)

						var _x = p.expression(0)

						localctx.(*ExpressionContext).onTrue = _x
					}

				}
				{
					p.SetState(548)
					p.Match(FqlParserColon)
				}
				{
					p.SetState(549)

					var _x = p.expression(3)

					localctx.(*ExpressionContext).onFalse = _x
				}

			}

		}
		p.SetState(554)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 55, p.GetParserRuleContext())
	}

	return localctx
}

// IPredicateContext is an interface to support dynamic dispatch.
type IPredicateContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetLeft returns the left rule contexts.
	GetLeft() IPredicateContext

	// GetRight returns the right rule contexts.
	GetRight() IPredicateContext

	// SetLeft sets the left rule contexts.
	SetLeft(IPredicateContext)

	// SetRight sets the right rule contexts.
	SetRight(IPredicateContext)

	// IsPredicateContext differentiates from other interfaces.
	IsPredicateContext()
}

type PredicateContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
	left   IPredicateContext
	right  IPredicateContext
}

func NewEmptyPredicateContext() *PredicateContext {
	var p = new(PredicateContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = FqlParserRULE_predicate
	return p
}

func (*PredicateContext) IsPredicateContext() {}

func NewPredicateContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *PredicateContext {
	var p = new(PredicateContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_predicate

	return p
}

func (s *PredicateContext) GetParser() antlr.Parser { return s.parser }

func (s *PredicateContext) GetLeft() IPredicateContext { return s.left }

func (s *PredicateContext) GetRight() IPredicateContext { return s.right }

func (s *PredicateContext) SetLeft(v IPredicateContext) { s.left = v }

func (s *PredicateContext) SetRight(v IPredicateContext) { s.right = v }

func (s *PredicateContext) ExpressionAtom() IExpressionAtomContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExpressionAtomContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IExpressionAtomContext)
}

func (s *PredicateContext) EqualityOperator() IEqualityOperatorContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IEqualityOperatorContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IEqualityOperatorContext)
}

func (s *PredicateContext) AllPredicate() []IPredicateContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IPredicateContext)(nil)).Elem())
	var tst = make([]IPredicateContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IPredicateContext)
		}
	}

	return tst
}

func (s *PredicateContext) Predicate(i int) IPredicateContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IPredicateContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IPredicateContext)
}

func (s *PredicateContext) ArrayOperator() IArrayOperatorContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IArrayOperatorContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IArrayOperatorContext)
}

func (s *PredicateContext) InOperator() IInOperatorContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IInOperatorContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IInOperatorContext)
}

func (s *PredicateContext) LikeOperator() ILikeOperatorContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ILikeOperatorContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ILikeOperatorContext)
}

func (s *PredicateContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *PredicateContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *PredicateContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterPredicate(s)
	}
}

func (s *PredicateContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitPredicate(s)
	}
}

func (s *PredicateContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitPredicate(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) Predicate() (localctx IPredicateContext) {
	return p.predicate(0)
}

func (p *FqlParser) predicate(_p int) (localctx IPredicateContext) {
	this := p
	_ = this

	var _parentctx antlr.ParserRuleContext = p.GetParserRuleContext()
	_parentState := p.GetState()
	localctx = NewPredicateContext(p, p.GetParserRuleContext(), _parentState)
	var _prevctx IPredicateContext = localctx
	var _ antlr.ParserRuleContext = _prevctx // TODO: To prevent unused variable warning.
	_startState := 118
	p.EnterRecursionRule(localctx, 118, FqlParserRULE_predicate, _p)

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
	{
		p.SetState(556)
		p.expressionAtom(0)
	}

	p.GetParserRuleContext().SetStop(p.GetTokenStream().LT(-1))
	p.SetState(576)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 57, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			if p.GetParseListeners() != nil {
				p.TriggerExitRuleEvent()
			}
			_prevctx = localctx
			p.SetState(574)
			p.GetErrorHandler().Sync(p)
			switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 56, p.GetParserRuleContext()) {
			case 1:
				localctx = NewPredicateContext(p, _parentctx, _parentState)
				localctx.(*PredicateContext).left = _prevctx
				p.PushNewRecursionContext(localctx, _startState, FqlParserRULE_predicate)
				p.SetState(558)

				if !(p.Precpred(p.GetParserRuleContext(), 5)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 5)", ""))
				}
				{
					p.SetState(559)
					p.EqualityOperator()
				}
				{
					p.SetState(560)

					var _x = p.predicate(6)

					localctx.(*PredicateContext).right = _x
				}

			case 2:
				localctx = NewPredicateContext(p, _parentctx, _parentState)
				localctx.(*PredicateContext).left = _prevctx
				p.PushNewRecursionContext(localctx, _startState, FqlParserRULE_predicate)
				p.SetState(562)

				if !(p.Precpred(p.GetParserRuleContext(), 4)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 4)", ""))
				}
				{
					p.SetState(563)
					p.ArrayOperator()
				}
				{
					p.SetState(564)

					var _x = p.predicate(5)

					localctx.(*PredicateContext).right = _x
				}

			case 3:
				localctx = NewPredicateContext(p, _parentctx, _parentState)
				localctx.(*PredicateContext).left = _prevctx
				p.PushNewRecursionContext(localctx, _startState, FqlParserRULE_predicate)
				p.SetState(566)

				if !(p.Precpred(p.GetParserRuleContext(), 3)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 3)", ""))
				}
				{
					p.SetState(567)
					p.InOperator()
				}
				{
					p.SetState(568)

					var _x = p.predicate(4)

					localctx.(*PredicateContext).right = _x
				}

			case 4:
				localctx = NewPredicateContext(p, _parentctx, _parentState)
				localctx.(*PredicateContext).left = _prevctx
				p.PushNewRecursionContext(localctx, _startState, FqlParserRULE_predicate)
				p.SetState(570)

				if !(p.Precpred(p.GetParserRuleContext(), 2)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 2)", ""))
				}
				{
					p.SetState(571)
					p.LikeOperator()
				}
				{
					p.SetState(572)

					var _x = p.predicate(3)

					localctx.(*PredicateContext).right = _x
				}

			}

		}
		p.SetState(578)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 57, p.GetParserRuleContext())
	}

	return localctx
}

// IExpressionAtomContext is an interface to support dynamic dispatch.
type IExpressionAtomContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetLeft returns the left rule contexts.
	GetLeft() IExpressionAtomContext

	// GetRight returns the right rule contexts.
	GetRight() IExpressionAtomContext

	// SetLeft sets the left rule contexts.
	SetLeft(IExpressionAtomContext)

	// SetRight sets the right rule contexts.
	SetRight(IExpressionAtomContext)

	// IsExpressionAtomContext differentiates from other interfaces.
	IsExpressionAtomContext()
}

type ExpressionAtomContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
	left   IExpressionAtomContext
	right  IExpressionAtomContext
}

func NewEmptyExpressionAtomContext() *ExpressionAtomContext {
	var p = new(ExpressionAtomContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = FqlParserRULE_expressionAtom
	return p
}

func (*ExpressionAtomContext) IsExpressionAtomContext() {}

func NewExpressionAtomContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ExpressionAtomContext {
	var p = new(ExpressionAtomContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_expressionAtom

	return p
}

func (s *ExpressionAtomContext) GetParser() antlr.Parser { return s.parser }

func (s *ExpressionAtomContext) GetLeft() IExpressionAtomContext { return s.left }

func (s *ExpressionAtomContext) GetRight() IExpressionAtomContext { return s.right }

func (s *ExpressionAtomContext) SetLeft(v IExpressionAtomContext) { s.left = v }

func (s *ExpressionAtomContext) SetRight(v IExpressionAtomContext) { s.right = v }

func (s *ExpressionAtomContext) FunctionCallExpression() IFunctionCallExpressionContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IFunctionCallExpressionContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IFunctionCallExpressionContext)
}

func (s *ExpressionAtomContext) RangeOperator() IRangeOperatorContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IRangeOperatorContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IRangeOperatorContext)
}

func (s *ExpressionAtomContext) Literal() ILiteralContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ILiteralContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ILiteralContext)
}

func (s *ExpressionAtomContext) Variable() IVariableContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IVariableContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IVariableContext)
}

func (s *ExpressionAtomContext) MemberExpression() IMemberExpressionContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IMemberExpressionContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IMemberExpressionContext)
}

func (s *ExpressionAtomContext) Param() IParamContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IParamContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IParamContext)
}

func (s *ExpressionAtomContext) OpenParen() antlr.TerminalNode {
	return s.GetToken(FqlParserOpenParen, 0)
}

func (s *ExpressionAtomContext) CloseParen() antlr.TerminalNode {
	return s.GetToken(FqlParserCloseParen, 0)
}

func (s *ExpressionAtomContext) ForExpression() IForExpressionContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IForExpressionContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IForExpressionContext)
}

func (s *ExpressionAtomContext) WaitForExpression() IWaitForExpressionContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IWaitForExpressionContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IWaitForExpressionContext)
}

func (s *ExpressionAtomContext) Expression() IExpressionContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExpressionContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *ExpressionAtomContext) ErrorOperator() IErrorOperatorContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IErrorOperatorContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IErrorOperatorContext)
}

func (s *ExpressionAtomContext) MultiplicativeOperator() IMultiplicativeOperatorContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IMultiplicativeOperatorContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IMultiplicativeOperatorContext)
}

func (s *ExpressionAtomContext) AllExpressionAtom() []IExpressionAtomContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IExpressionAtomContext)(nil)).Elem())
	var tst = make([]IExpressionAtomContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IExpressionAtomContext)
		}
	}

	return tst
}

func (s *ExpressionAtomContext) ExpressionAtom(i int) IExpressionAtomContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExpressionAtomContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IExpressionAtomContext)
}

func (s *ExpressionAtomContext) AdditiveOperator() IAdditiveOperatorContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IAdditiveOperatorContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IAdditiveOperatorContext)
}

func (s *ExpressionAtomContext) RegexpOperator() IRegexpOperatorContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IRegexpOperatorContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IRegexpOperatorContext)
}

func (s *ExpressionAtomContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ExpressionAtomContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ExpressionAtomContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterExpressionAtom(s)
	}
}

func (s *ExpressionAtomContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitExpressionAtom(s)
	}
}

func (s *ExpressionAtomContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitExpressionAtom(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) ExpressionAtom() (localctx IExpressionAtomContext) {
	return p.expressionAtom(0)
}

func (p *FqlParser) expressionAtom(_p int) (localctx IExpressionAtomContext) {
	this := p
	_ = this

	var _parentctx antlr.ParserRuleContext = p.GetParserRuleContext()
	_parentState := p.GetState()
	localctx = NewExpressionAtomContext(p, p.GetParserRuleContext(), _parentState)
	var _prevctx IExpressionAtomContext = localctx
	var _ antlr.ParserRuleContext = _prevctx // TODO: To prevent unused variable warning.
	_startState := 120
	p.EnterRecursionRule(localctx, 120, FqlParserRULE_expressionAtom, _p)

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
	p.SetState(596)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 60, p.GetParserRuleContext()) {
	case 1:
		{
			p.SetState(580)
			p.FunctionCallExpression()
		}

	case 2:
		{
			p.SetState(581)
			p.RangeOperator()
		}

	case 3:
		{
			p.SetState(582)
			p.Literal()
		}

	case 4:
		{
			p.SetState(583)
			p.Variable()
		}

	case 5:
		{
			p.SetState(584)
			p.MemberExpression()
		}

	case 6:
		{
			p.SetState(585)
			p.Param()
		}

	case 7:
		{
			p.SetState(586)
			p.Match(FqlParserOpenParen)
		}
		p.SetState(590)
		p.GetErrorHandler().Sync(p)
		switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 58, p.GetParserRuleContext()) {
		case 1:
			{
				p.SetState(587)
				p.ForExpression()
			}

		case 2:
			{
				p.SetState(588)
				p.WaitForExpression()
			}

		case 3:
			{
				p.SetState(589)
				p.expression(0)
			}

		}
		{
			p.SetState(592)
			p.Match(FqlParserCloseParen)
		}
		p.SetState(594)
		p.GetErrorHandler().Sync(p)

		if p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 59, p.GetParserRuleContext()) == 1 {
			{
				p.SetState(593)
				p.ErrorOperator()
			}

		}

	}
	p.GetParserRuleContext().SetStop(p.GetTokenStream().LT(-1))
	p.SetState(612)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 62, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			if p.GetParseListeners() != nil {
				p.TriggerExitRuleEvent()
			}
			_prevctx = localctx
			p.SetState(610)
			p.GetErrorHandler().Sync(p)
			switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 61, p.GetParserRuleContext()) {
			case 1:
				localctx = NewExpressionAtomContext(p, _parentctx, _parentState)
				localctx.(*ExpressionAtomContext).left = _prevctx
				p.PushNewRecursionContext(localctx, _startState, FqlParserRULE_expressionAtom)
				p.SetState(598)

				if !(p.Precpred(p.GetParserRuleContext(), 10)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 10)", ""))
				}
				{
					p.SetState(599)
					p.MultiplicativeOperator()
				}
				{
					p.SetState(600)

					var _x = p.expressionAtom(11)

					localctx.(*ExpressionAtomContext).right = _x
				}

			case 2:
				localctx = NewExpressionAtomContext(p, _parentctx, _parentState)
				localctx.(*ExpressionAtomContext).left = _prevctx
				p.PushNewRecursionContext(localctx, _startState, FqlParserRULE_expressionAtom)
				p.SetState(602)

				if !(p.Precpred(p.GetParserRuleContext(), 9)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 9)", ""))
				}
				{
					p.SetState(603)
					p.AdditiveOperator()
				}
				{
					p.SetState(604)

					var _x = p.expressionAtom(10)

					localctx.(*ExpressionAtomContext).right = _x
				}

			case 3:
				localctx = NewExpressionAtomContext(p, _parentctx, _parentState)
				localctx.(*ExpressionAtomContext).left = _prevctx
				p.PushNewRecursionContext(localctx, _startState, FqlParserRULE_expressionAtom)
				p.SetState(606)

				if !(p.Precpred(p.GetParserRuleContext(), 8)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 8)", ""))
				}
				{
					p.SetState(607)
					p.RegexpOperator()
				}
				{
					p.SetState(608)

					var _x = p.expressionAtom(9)

					localctx.(*ExpressionAtomContext).right = _x
				}

			}

		}
		p.SetState(614)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 62, p.GetParserRuleContext())
	}

	return localctx
}

// IArrayOperatorContext is an interface to support dynamic dispatch.
type IArrayOperatorContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetOperator returns the operator token.
	GetOperator() antlr.Token

	// SetOperator sets the operator token.
	SetOperator(antlr.Token)

	// IsArrayOperatorContext differentiates from other interfaces.
	IsArrayOperatorContext()
}

type ArrayOperatorContext struct {
	*antlr.BaseParserRuleContext
	parser   antlr.Parser
	operator antlr.Token
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

func (s *ArrayOperatorContext) GetOperator() antlr.Token { return s.operator }

func (s *ArrayOperatorContext) SetOperator(v antlr.Token) { s.operator = v }

func (s *ArrayOperatorContext) All() antlr.TerminalNode {
	return s.GetToken(FqlParserAll, 0)
}

func (s *ArrayOperatorContext) Any() antlr.TerminalNode {
	return s.GetToken(FqlParserAny, 0)
}

func (s *ArrayOperatorContext) None() antlr.TerminalNode {
	return s.GetToken(FqlParserNone, 0)
}

func (s *ArrayOperatorContext) InOperator() IInOperatorContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IInOperatorContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IInOperatorContext)
}

func (s *ArrayOperatorContext) EqualityOperator() IEqualityOperatorContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IEqualityOperatorContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IEqualityOperatorContext)
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
	this := p
	_ = this

	localctx = NewArrayOperatorContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 122, FqlParserRULE_arrayOperator)
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
		p.SetState(615)

		var _lt = p.GetTokenStream().LT(1)

		localctx.(*ArrayOperatorContext).operator = _lt

		_la = p.GetTokenStream().LA(1)

		if !(((_la-48)&-(0x1f+1)) == 0 && ((1<<uint((_la-48)))&((1<<(FqlParserNone-48))|(1<<(FqlParserAll-48))|(1<<(FqlParserAny-48)))) != 0) {
			var _ri = p.GetErrorHandler().RecoverInline(p)

			localctx.(*ArrayOperatorContext).operator = _ri
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}
	p.SetState(618)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case FqlParserNot, FqlParserIn:
		{
			p.SetState(616)
			p.InOperator()
		}

	case FqlParserGt, FqlParserLt, FqlParserEq, FqlParserGte, FqlParserLte, FqlParserNeq:
		{
			p.SetState(617)
			p.EqualityOperator()
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
	this := p
	_ = this

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
		p.SetState(620)
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
	this := p
	_ = this

	localctx = NewInOperatorContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 126, FqlParserRULE_inOperator)
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
	p.SetState(623)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == FqlParserNot {
		{
			p.SetState(622)
			p.Match(FqlParserNot)
		}

	}
	{
		p.SetState(625)
		p.Match(FqlParserIn)
	}

	return localctx
}

// ILikeOperatorContext is an interface to support dynamic dispatch.
type ILikeOperatorContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsLikeOperatorContext differentiates from other interfaces.
	IsLikeOperatorContext()
}

type LikeOperatorContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyLikeOperatorContext() *LikeOperatorContext {
	var p = new(LikeOperatorContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = FqlParserRULE_likeOperator
	return p
}

func (*LikeOperatorContext) IsLikeOperatorContext() {}

func NewLikeOperatorContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *LikeOperatorContext {
	var p = new(LikeOperatorContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_likeOperator

	return p
}

func (s *LikeOperatorContext) GetParser() antlr.Parser { return s.parser }

func (s *LikeOperatorContext) Like() antlr.TerminalNode {
	return s.GetToken(FqlParserLike, 0)
}

func (s *LikeOperatorContext) Not() antlr.TerminalNode {
	return s.GetToken(FqlParserNot, 0)
}

func (s *LikeOperatorContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *LikeOperatorContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *LikeOperatorContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterLikeOperator(s)
	}
}

func (s *LikeOperatorContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitLikeOperator(s)
	}
}

func (s *LikeOperatorContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitLikeOperator(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) LikeOperator() (localctx ILikeOperatorContext) {
	this := p
	_ = this

	localctx = NewLikeOperatorContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 128, FqlParserRULE_likeOperator)
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
	p.SetState(628)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == FqlParserNot {
		{
			p.SetState(627)
			p.Match(FqlParserNot)
		}

	}
	{
		p.SetState(630)
		p.Match(FqlParserLike)
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
	this := p
	_ = this

	localctx = NewUnaryOperatorContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 130, FqlParserRULE_unaryOperator)
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
		p.SetState(632)
		_la = p.GetTokenStream().LA(1)

		if !(_la == FqlParserPlus || _la == FqlParserMinus || _la == FqlParserNot) {
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
	this := p
	_ = this

	localctx = NewRegexpOperatorContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 132, FqlParserRULE_regexpOperator)
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
		p.SetState(634)
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
	this := p
	_ = this

	localctx = NewLogicalAndOperatorContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 134, FqlParserRULE_logicalAndOperator)

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
		p.SetState(636)
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
	this := p
	_ = this

	localctx = NewLogicalOrOperatorContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 136, FqlParserRULE_logicalOrOperator)

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
		p.SetState(638)
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
	this := p
	_ = this

	localctx = NewMultiplicativeOperatorContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 138, FqlParserRULE_multiplicativeOperator)
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
		p.SetState(640)
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
	this := p
	_ = this

	localctx = NewAdditiveOperatorContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 140, FqlParserRULE_additiveOperator)
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

		if !(_la == FqlParserPlus || _la == FqlParserMinus) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}

	return localctx
}

// IErrorOperatorContext is an interface to support dynamic dispatch.
type IErrorOperatorContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsErrorOperatorContext differentiates from other interfaces.
	IsErrorOperatorContext()
}

type ErrorOperatorContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyErrorOperatorContext() *ErrorOperatorContext {
	var p = new(ErrorOperatorContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = FqlParserRULE_errorOperator
	return p
}

func (*ErrorOperatorContext) IsErrorOperatorContext() {}

func NewErrorOperatorContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ErrorOperatorContext {
	var p = new(ErrorOperatorContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_errorOperator

	return p
}

func (s *ErrorOperatorContext) GetParser() antlr.Parser { return s.parser }

func (s *ErrorOperatorContext) QuestionMark() antlr.TerminalNode {
	return s.GetToken(FqlParserQuestionMark, 0)
}

func (s *ErrorOperatorContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ErrorOperatorContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ErrorOperatorContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterErrorOperator(s)
	}
}

func (s *ErrorOperatorContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitErrorOperator(s)
	}
}

func (s *ErrorOperatorContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitErrorOperator(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) ErrorOperator() (localctx IErrorOperatorContext) {
	this := p
	_ = this

	localctx = NewErrorOperatorContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 142, FqlParserRULE_errorOperator)

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
		p.SetState(644)
		p.Match(FqlParserQuestionMark)
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

	case 59:
		var t *PredicateContext = nil
		if localctx != nil {
			t = localctx.(*PredicateContext)
		}
		return p.Predicate_Sempred(t, predIndex)

	case 60:
		var t *ExpressionAtomContext = nil
		if localctx != nil {
			t = localctx.(*ExpressionAtomContext)
		}
		return p.ExpressionAtom_Sempred(t, predIndex)

	default:
		panic("No predicate with index: " + fmt.Sprint(ruleIndex))
	}
}

func (p *FqlParser) Expression_Sempred(localctx antlr.RuleContext, predIndex int) bool {
	this := p
	_ = this

	switch predIndex {
	case 0:
		return p.Precpred(p.GetParserRuleContext(), 4)

	case 1:
		return p.Precpred(p.GetParserRuleContext(), 3)

	case 2:
		return p.Precpred(p.GetParserRuleContext(), 2)

	default:
		panic("No predicate with index: " + fmt.Sprint(predIndex))
	}
}

func (p *FqlParser) Predicate_Sempred(localctx antlr.RuleContext, predIndex int) bool {
	this := p
	_ = this

	switch predIndex {
	case 3:
		return p.Precpred(p.GetParserRuleContext(), 5)

	case 4:
		return p.Precpred(p.GetParserRuleContext(), 4)

	case 5:
		return p.Precpred(p.GetParserRuleContext(), 3)

	case 6:
		return p.Precpred(p.GetParserRuleContext(), 2)

	default:
		panic("No predicate with index: " + fmt.Sprint(predIndex))
	}
}

func (p *FqlParser) ExpressionAtom_Sempred(localctx antlr.RuleContext, predIndex int) bool {
	this := p
	_ = this

	switch predIndex {
	case 7:
		return p.Precpred(p.GetParserRuleContext(), 10)

	case 8:
		return p.Precpred(p.GetParserRuleContext(), 9)

	case 9:
		return p.Precpred(p.GetParserRuleContext(), 8)

	default:
		panic("No predicate with index: " + fmt.Sprint(predIndex))
	}
}
