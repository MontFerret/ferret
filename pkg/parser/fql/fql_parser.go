// Code generated from antlr/FqlParser.g4 by ANTLR 4.9.2. DO NOT EDIT.

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
	3, 24715, 42794, 33075, 47597, 16764, 15335, 30598, 22884, 3, 71, 616,
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
	2, 7, 2, 142, 10, 2, 12, 2, 14, 2, 145, 11, 2, 3, 2, 3, 2, 3, 2, 3, 3,
	3, 3, 3, 4, 3, 4, 3, 5, 3, 5, 3, 5, 3, 6, 3, 6, 3, 6, 3, 7, 7, 7, 161,
	10, 7, 12, 7, 14, 7, 164, 11, 7, 3, 7, 3, 7, 3, 8, 3, 8, 3, 8, 5, 8, 171,
	10, 8, 3, 9, 3, 9, 5, 9, 175, 10, 9, 3, 10, 3, 10, 3, 10, 3, 10, 3, 10,
	3, 11, 3, 11, 5, 11, 184, 10, 11, 3, 11, 3, 11, 3, 12, 3, 12, 3, 12, 3,
	12, 5, 12, 192, 10, 12, 3, 13, 3, 13, 5, 13, 196, 10, 13, 3, 14, 3, 14,
	3, 14, 3, 14, 5, 14, 202, 10, 14, 3, 14, 3, 14, 3, 14, 7, 14, 207, 10,
	14, 12, 14, 14, 14, 210, 11, 14, 3, 14, 3, 14, 3, 14, 3, 14, 3, 14, 5,
	14, 217, 10, 14, 3, 14, 3, 14, 3, 14, 7, 14, 222, 10, 14, 12, 14, 14, 14,
	225, 11, 14, 3, 14, 3, 14, 5, 14, 229, 10, 14, 3, 15, 3, 15, 3, 15, 3,
	15, 3, 15, 3, 15, 3, 15, 5, 15, 238, 10, 15, 3, 16, 3, 16, 3, 16, 3, 16,
	5, 16, 244, 10, 16, 3, 17, 3, 17, 5, 17, 248, 10, 17, 3, 18, 3, 18, 5,
	18, 252, 10, 18, 3, 19, 3, 19, 5, 19, 256, 10, 19, 3, 20, 3, 20, 3, 20,
	3, 21, 3, 21, 3, 21, 3, 21, 5, 21, 265, 10, 21, 3, 22, 3, 22, 5, 22, 269,
	10, 22, 3, 23, 3, 23, 3, 23, 3, 23, 7, 23, 275, 10, 23, 12, 23, 14, 23,
	278, 11, 23, 3, 24, 3, 24, 5, 24, 282, 10, 24, 3, 25, 3, 25, 3, 25, 3,
	25, 3, 25, 3, 25, 3, 25, 3, 25, 3, 25, 3, 25, 3, 25, 3, 25, 3, 25, 3, 25,
	3, 25, 3, 25, 3, 25, 3, 25, 5, 25, 302, 10, 25, 3, 26, 3, 26, 3, 26, 3,
	26, 3, 27, 3, 27, 3, 27, 7, 27, 311, 10, 27, 12, 27, 14, 27, 314, 11, 27,
	3, 28, 3, 28, 3, 28, 3, 28, 7, 28, 320, 10, 28, 12, 28, 14, 28, 323, 11,
	28, 3, 29, 3, 29, 3, 29, 3, 29, 3, 30, 3, 30, 3, 30, 3, 30, 3, 30, 3, 30,
	5, 30, 335, 10, 30, 5, 30, 337, 10, 30, 3, 31, 3, 31, 3, 31, 3, 31, 3,
	31, 3, 32, 3, 32, 3, 32, 3, 32, 5, 32, 348, 10, 32, 3, 33, 3, 33, 3, 33,
	3, 33, 3, 33, 3, 33, 5, 33, 356, 10, 33, 3, 33, 5, 33, 359, 10, 33, 3,
	34, 3, 34, 3, 34, 5, 34, 364, 10, 34, 3, 35, 3, 35, 3, 35, 3, 35, 3, 35,
	5, 35, 371, 10, 35, 3, 36, 3, 36, 3, 36, 5, 36, 376, 10, 36, 3, 37, 3,
	37, 3, 37, 5, 37, 381, 10, 37, 3, 37, 3, 37, 3, 37, 3, 37, 5, 37, 387,
	10, 37, 3, 38, 3, 38, 3, 38, 3, 38, 7, 38, 393, 10, 38, 12, 38, 14, 38,
	396, 11, 38, 3, 38, 5, 38, 399, 10, 38, 5, 38, 401, 10, 38, 3, 38, 3, 38,
	3, 39, 3, 39, 3, 39, 3, 39, 7, 39, 409, 10, 39, 12, 39, 14, 39, 412, 11,
	39, 3, 39, 5, 39, 415, 10, 39, 5, 39, 417, 10, 39, 3, 39, 3, 39, 3, 40,
	3, 40, 3, 40, 3, 40, 3, 40, 3, 40, 3, 40, 3, 40, 3, 40, 5, 40, 430, 10,
	40, 3, 41, 3, 41, 3, 41, 3, 41, 3, 42, 3, 42, 3, 42, 5, 42, 439, 10, 42,
	3, 43, 3, 43, 3, 44, 3, 44, 3, 45, 3, 45, 3, 46, 3, 46, 3, 47, 3, 47, 3,
	48, 3, 48, 3, 48, 3, 48, 3, 48, 3, 48, 3, 48, 3, 48, 3, 48, 3, 48, 3, 48,
	3, 48, 3, 48, 3, 48, 3, 48, 3, 48, 3, 48, 3, 48, 5, 48, 469, 10, 48, 3,
	48, 3, 48, 3, 48, 3, 48, 3, 48, 3, 48, 3, 48, 3, 48, 3, 48, 3, 48, 3, 48,
	3, 48, 5, 48, 483, 10, 48, 3, 48, 3, 48, 3, 48, 3, 48, 3, 48, 3, 48, 3,
	48, 3, 48, 3, 48, 3, 48, 3, 48, 3, 48, 3, 48, 3, 48, 3, 48, 3, 48, 3, 48,
	3, 48, 3, 48, 3, 48, 3, 48, 3, 48, 3, 48, 3, 48, 3, 48, 3, 48, 3, 48, 3,
	48, 3, 48, 5, 48, 514, 10, 48, 3, 48, 3, 48, 7, 48, 518, 10, 48, 12, 48,
	14, 48, 521, 11, 48, 3, 49, 3, 49, 3, 49, 3, 49, 5, 49, 527, 10, 49, 3,
	50, 3, 50, 6, 50, 531, 10, 50, 13, 50, 14, 50, 532, 3, 51, 3, 51, 3, 51,
	3, 51, 3, 51, 5, 51, 540, 10, 51, 3, 52, 3, 52, 3, 52, 3, 52, 3, 53, 3,
	53, 5, 53, 548, 10, 53, 3, 54, 5, 54, 551, 10, 54, 3, 54, 3, 54, 3, 54,
	3, 54, 5, 54, 557, 10, 54, 3, 54, 5, 54, 560, 10, 54, 3, 55, 3, 55, 3,
	56, 3, 56, 3, 57, 7, 57, 567, 10, 57, 12, 57, 14, 57, 570, 11, 57, 3, 58,
	3, 58, 3, 58, 3, 58, 7, 58, 576, 10, 58, 12, 58, 14, 58, 579, 11, 58, 5,
	58, 581, 10, 58, 3, 58, 3, 58, 3, 59, 3, 59, 3, 60, 3, 60, 3, 60, 5, 60,
	590, 10, 60, 3, 61, 3, 61, 3, 61, 5, 61, 595, 10, 61, 3, 62, 3, 62, 3,
	63, 3, 63, 3, 64, 3, 64, 3, 65, 3, 65, 3, 66, 3, 66, 3, 67, 3, 67, 3, 68,
	3, 68, 3, 69, 3, 69, 3, 69, 3, 70, 3, 70, 3, 70, 2, 3, 94, 71, 2, 4, 6,
	8, 10, 12, 14, 16, 18, 20, 22, 24, 26, 28, 30, 32, 34, 36, 38, 40, 42,
	44, 46, 48, 50, 52, 54, 56, 58, 60, 62, 64, 66, 68, 70, 72, 74, 76, 78,
	80, 82, 84, 86, 88, 90, 92, 94, 96, 98, 100, 102, 104, 106, 108, 110, 112,
	114, 116, 118, 120, 122, 124, 126, 128, 130, 132, 134, 136, 138, 2, 10,
	3, 2, 48, 49, 6, 2, 30, 31, 37, 39, 41, 62, 66, 66, 4, 2, 48, 48, 56, 57,
	3, 2, 17, 22, 3, 2, 35, 36, 3, 2, 23, 25, 3, 2, 26, 27, 4, 2, 26, 27, 61,
	61, 2, 650, 2, 143, 3, 2, 2, 2, 4, 149, 3, 2, 2, 2, 6, 151, 3, 2, 2, 2,
	8, 153, 3, 2, 2, 2, 10, 156, 3, 2, 2, 2, 12, 162, 3, 2, 2, 2, 14, 170,
	3, 2, 2, 2, 16, 174, 3, 2, 2, 2, 18, 176, 3, 2, 2, 2, 20, 181, 3, 2, 2,
	2, 22, 187, 3, 2, 2, 2, 24, 195, 3, 2, 2, 2, 26, 228, 3, 2, 2, 2, 28, 237,
	3, 2, 2, 2, 30, 243, 3, 2, 2, 2, 32, 247, 3, 2, 2, 2, 34, 251, 3, 2, 2,
	2, 36, 255, 3, 2, 2, 2, 38, 257, 3, 2, 2, 2, 40, 260, 3, 2, 2, 2, 42, 268,
	3, 2, 2, 2, 44, 270, 3, 2, 2, 2, 46, 279, 3, 2, 2, 2, 48, 301, 3, 2, 2,
	2, 50, 303, 3, 2, 2, 2, 52, 307, 3, 2, 2, 2, 54, 315, 3, 2, 2, 2, 56, 324,
	3, 2, 2, 2, 58, 336, 3, 2, 2, 2, 60, 338, 3, 2, 2, 2, 62, 347, 3, 2, 2,
	2, 64, 349, 3, 2, 2, 2, 66, 363, 3, 2, 2, 2, 68, 370, 3, 2, 2, 2, 70, 375,
	3, 2, 2, 2, 72, 380, 3, 2, 2, 2, 74, 388, 3, 2, 2, 2, 76, 404, 3, 2, 2,
	2, 78, 429, 3, 2, 2, 2, 80, 431, 3, 2, 2, 2, 82, 438, 3, 2, 2, 2, 84, 440,
	3, 2, 2, 2, 86, 442, 3, 2, 2, 2, 88, 444, 3, 2, 2, 2, 90, 446, 3, 2, 2,
	2, 92, 448, 3, 2, 2, 2, 94, 468, 3, 2, 2, 2, 96, 522, 3, 2, 2, 2, 98, 528,
	3, 2, 2, 2, 100, 539, 3, 2, 2, 2, 102, 541, 3, 2, 2, 2, 104, 545, 3, 2,
	2, 2, 106, 559, 3, 2, 2, 2, 108, 561, 3, 2, 2, 2, 110, 563, 3, 2, 2, 2,
	112, 568, 3, 2, 2, 2, 114, 571, 3, 2, 2, 2, 116, 584, 3, 2, 2, 2, 118,
	589, 3, 2, 2, 2, 120, 594, 3, 2, 2, 2, 122, 596, 3, 2, 2, 2, 124, 598,
	3, 2, 2, 2, 126, 600, 3, 2, 2, 2, 128, 602, 3, 2, 2, 2, 130, 604, 3, 2,
	2, 2, 132, 606, 3, 2, 2, 2, 134, 608, 3, 2, 2, 2, 136, 610, 3, 2, 2, 2,
	138, 613, 3, 2, 2, 2, 140, 142, 5, 4, 3, 2, 141, 140, 3, 2, 2, 2, 142,
	145, 3, 2, 2, 2, 143, 141, 3, 2, 2, 2, 143, 144, 3, 2, 2, 2, 144, 146,
	3, 2, 2, 2, 145, 143, 3, 2, 2, 2, 146, 147, 5, 12, 7, 2, 147, 148, 7, 2,
	2, 3, 148, 3, 3, 2, 2, 2, 149, 150, 5, 6, 4, 2, 150, 5, 3, 2, 2, 2, 151,
	152, 5, 8, 5, 2, 152, 7, 3, 2, 2, 2, 153, 154, 7, 51, 2, 2, 154, 155, 5,
	10, 6, 2, 155, 9, 3, 2, 2, 2, 156, 157, 5, 112, 57, 2, 157, 158, 7, 66,
	2, 2, 158, 11, 3, 2, 2, 2, 159, 161, 5, 14, 8, 2, 160, 159, 3, 2, 2, 2,
	161, 164, 3, 2, 2, 2, 162, 160, 3, 2, 2, 2, 162, 163, 3, 2, 2, 2, 163,
	165, 3, 2, 2, 2, 164, 162, 3, 2, 2, 2, 165, 166, 5, 16, 9, 2, 166, 13,
	3, 2, 2, 2, 167, 171, 5, 18, 10, 2, 168, 171, 5, 104, 53, 2, 169, 171,
	5, 64, 33, 2, 170, 167, 3, 2, 2, 2, 170, 168, 3, 2, 2, 2, 170, 169, 3,
	2, 2, 2, 171, 15, 3, 2, 2, 2, 172, 175, 5, 20, 11, 2, 173, 175, 5, 26,
	14, 2, 174, 172, 3, 2, 2, 2, 174, 173, 3, 2, 2, 2, 175, 17, 3, 2, 2, 2,
	176, 177, 7, 45, 2, 2, 177, 178, 7, 66, 2, 2, 178, 179, 7, 33, 2, 2, 179,
	180, 5, 94, 48, 2, 180, 19, 3, 2, 2, 2, 181, 183, 7, 38, 2, 2, 182, 184,
	7, 41, 2, 2, 183, 182, 3, 2, 2, 2, 183, 184, 3, 2, 2, 2, 184, 185, 3, 2,
	2, 2, 185, 186, 5, 94, 48, 2, 186, 21, 3, 2, 2, 2, 187, 188, 7, 13, 2,
	2, 188, 189, 5, 24, 13, 2, 189, 191, 7, 14, 2, 2, 190, 192, 5, 108, 55,
	2, 191, 190, 3, 2, 2, 2, 191, 192, 3, 2, 2, 2, 192, 23, 3, 2, 2, 2, 193,
	196, 5, 26, 14, 2, 194, 196, 5, 64, 33, 2, 195, 193, 3, 2, 2, 2, 195, 194,
	3, 2, 2, 2, 196, 25, 3, 2, 2, 2, 197, 198, 7, 37, 2, 2, 198, 201, 7, 66,
	2, 2, 199, 200, 7, 10, 2, 2, 200, 202, 7, 66, 2, 2, 201, 199, 3, 2, 2,
	2, 201, 202, 3, 2, 2, 2, 202, 203, 3, 2, 2, 2, 203, 204, 7, 62, 2, 2, 204,
	208, 5, 28, 15, 2, 205, 207, 5, 34, 18, 2, 206, 205, 3, 2, 2, 2, 207, 210,
	3, 2, 2, 2, 208, 206, 3, 2, 2, 2, 208, 209, 3, 2, 2, 2, 209, 211, 3, 2,
	2, 2, 210, 208, 3, 2, 2, 2, 211, 212, 5, 36, 19, 2, 212, 229, 3, 2, 2,
	2, 213, 214, 7, 37, 2, 2, 214, 216, 7, 66, 2, 2, 215, 217, 7, 63, 2, 2,
	216, 215, 3, 2, 2, 2, 216, 217, 3, 2, 2, 2, 217, 218, 3, 2, 2, 2, 218,
	219, 7, 64, 2, 2, 219, 223, 5, 94, 48, 2, 220, 222, 5, 34, 18, 2, 221,
	220, 3, 2, 2, 2, 222, 225, 3, 2, 2, 2, 223, 221, 3, 2, 2, 2, 223, 224,
	3, 2, 2, 2, 224, 226, 3, 2, 2, 2, 225, 223, 3, 2, 2, 2, 226, 227, 5, 36,
	19, 2, 227, 229, 3, 2, 2, 2, 228, 197, 3, 2, 2, 2, 228, 213, 3, 2, 2, 2,
	229, 27, 3, 2, 2, 2, 230, 238, 5, 104, 53, 2, 231, 238, 5, 74, 38, 2, 232,
	238, 5, 76, 39, 2, 233, 238, 5, 138, 70, 2, 234, 238, 5, 98, 50, 2, 235,
	238, 5, 72, 37, 2, 236, 238, 5, 136, 69, 2, 237, 230, 3, 2, 2, 2, 237,
	231, 3, 2, 2, 2, 237, 232, 3, 2, 2, 2, 237, 233, 3, 2, 2, 2, 237, 234,
	3, 2, 2, 2, 237, 235, 3, 2, 2, 2, 237, 236, 3, 2, 2, 2, 238, 29, 3, 2,
	2, 2, 239, 244, 5, 40, 21, 2, 240, 244, 5, 44, 23, 2, 241, 244, 5, 38,
	20, 2, 242, 244, 5, 48, 25, 2, 243, 239, 3, 2, 2, 2, 243, 240, 3, 2, 2,
	2, 243, 241, 3, 2, 2, 2, 243, 242, 3, 2, 2, 2, 244, 31, 3, 2, 2, 2, 245,
	248, 5, 18, 10, 2, 246, 248, 5, 104, 53, 2, 247, 245, 3, 2, 2, 2, 247,
	246, 3, 2, 2, 2, 248, 33, 3, 2, 2, 2, 249, 252, 5, 32, 17, 2, 250, 252,
	5, 30, 16, 2, 251, 249, 3, 2, 2, 2, 251, 250, 3, 2, 2, 2, 252, 35, 3, 2,
	2, 2, 253, 256, 5, 20, 11, 2, 254, 256, 5, 26, 14, 2, 255, 253, 3, 2, 2,
	2, 255, 254, 3, 2, 2, 2, 256, 37, 3, 2, 2, 2, 257, 258, 7, 42, 2, 2, 258,
	259, 5, 94, 48, 2, 259, 39, 3, 2, 2, 2, 260, 261, 7, 44, 2, 2, 261, 264,
	5, 42, 22, 2, 262, 263, 7, 10, 2, 2, 263, 265, 5, 42, 22, 2, 264, 262,
	3, 2, 2, 2, 264, 265, 3, 2, 2, 2, 265, 41, 3, 2, 2, 2, 266, 269, 7, 68,
	2, 2, 267, 269, 5, 136, 69, 2, 268, 266, 3, 2, 2, 2, 268, 267, 3, 2, 2,
	2, 269, 43, 3, 2, 2, 2, 270, 271, 7, 43, 2, 2, 271, 276, 5, 46, 24, 2,
	272, 273, 7, 10, 2, 2, 273, 275, 5, 46, 24, 2, 274, 272, 3, 2, 2, 2, 275,
	278, 3, 2, 2, 2, 276, 274, 3, 2, 2, 2, 276, 277, 3, 2, 2, 2, 277, 45, 3,
	2, 2, 2, 278, 276, 3, 2, 2, 2, 279, 281, 5, 94, 48, 2, 280, 282, 7, 47,
	2, 2, 281, 280, 3, 2, 2, 2, 281, 282, 3, 2, 2, 2, 282, 47, 3, 2, 2, 2,
	283, 284, 7, 46, 2, 2, 284, 302, 5, 60, 31, 2, 285, 286, 7, 46, 2, 2, 286,
	302, 5, 54, 28, 2, 287, 288, 7, 46, 2, 2, 288, 289, 5, 52, 27, 2, 289,
	290, 5, 54, 28, 2, 290, 302, 3, 2, 2, 2, 291, 292, 7, 46, 2, 2, 292, 293,
	5, 52, 27, 2, 293, 294, 5, 58, 30, 2, 294, 302, 3, 2, 2, 2, 295, 296, 7,
	46, 2, 2, 296, 297, 5, 52, 27, 2, 297, 298, 5, 60, 31, 2, 298, 302, 3,
	2, 2, 2, 299, 300, 7, 46, 2, 2, 300, 302, 5, 52, 27, 2, 301, 283, 3, 2,
	2, 2, 301, 285, 3, 2, 2, 2, 301, 287, 3, 2, 2, 2, 301, 291, 3, 2, 2, 2,
	301, 295, 3, 2, 2, 2, 301, 299, 3, 2, 2, 2, 302, 49, 3, 2, 2, 2, 303, 304,
	7, 66, 2, 2, 304, 305, 7, 33, 2, 2, 305, 306, 5, 94, 48, 2, 306, 51, 3,
	2, 2, 2, 307, 312, 5, 50, 26, 2, 308, 309, 7, 10, 2, 2, 309, 311, 5, 50,
	26, 2, 310, 308, 3, 2, 2, 2, 311, 314, 3, 2, 2, 2, 312, 310, 3, 2, 2, 2,
	312, 313, 3, 2, 2, 2, 313, 53, 3, 2, 2, 2, 314, 312, 3, 2, 2, 2, 315, 316,
	7, 58, 2, 2, 316, 321, 5, 56, 29, 2, 317, 318, 7, 10, 2, 2, 318, 320, 5,
	56, 29, 2, 319, 317, 3, 2, 2, 2, 320, 323, 3, 2, 2, 2, 321, 319, 3, 2,
	2, 2, 321, 322, 3, 2, 2, 2, 322, 55, 3, 2, 2, 2, 323, 321, 3, 2, 2, 2,
	324, 325, 7, 66, 2, 2, 325, 326, 7, 33, 2, 2, 326, 327, 5, 104, 53, 2,
	327, 57, 3, 2, 2, 2, 328, 329, 7, 52, 2, 2, 329, 337, 5, 50, 26, 2, 330,
	331, 7, 52, 2, 2, 331, 334, 7, 66, 2, 2, 332, 333, 7, 53, 2, 2, 333, 335,
	7, 66, 2, 2, 334, 332, 3, 2, 2, 2, 334, 335, 3, 2, 2, 2, 335, 337, 3, 2,
	2, 2, 336, 328, 3, 2, 2, 2, 336, 330, 3, 2, 2, 2, 337, 59, 3, 2, 2, 2,
	338, 339, 7, 54, 2, 2, 339, 340, 7, 55, 2, 2, 340, 341, 7, 52, 2, 2, 341,
	342, 7, 66, 2, 2, 342, 61, 3, 2, 2, 2, 343, 344, 7, 40, 2, 2, 344, 348,
	5, 76, 39, 2, 345, 346, 7, 40, 2, 2, 346, 348, 5, 138, 70, 2, 347, 343,
	3, 2, 2, 2, 347, 345, 3, 2, 2, 2, 348, 63, 3, 2, 2, 2, 349, 350, 7, 39,
	2, 2, 350, 351, 7, 59, 2, 2, 351, 352, 5, 68, 35, 2, 352, 353, 7, 62, 2,
	2, 353, 355, 5, 70, 36, 2, 354, 356, 5, 62, 32, 2, 355, 354, 3, 2, 2, 2,
	355, 356, 3, 2, 2, 2, 356, 358, 3, 2, 2, 2, 357, 359, 5, 66, 34, 2, 358,
	357, 3, 2, 2, 2, 358, 359, 3, 2, 2, 2, 359, 65, 3, 2, 2, 2, 360, 364, 5,
	90, 46, 2, 361, 364, 5, 138, 70, 2, 362, 364, 5, 136, 69, 2, 363, 360,
	3, 2, 2, 2, 363, 361, 3, 2, 2, 2, 363, 362, 3, 2, 2, 2, 364, 67, 3, 2,
	2, 2, 365, 371, 5, 86, 44, 2, 366, 371, 5, 138, 70, 2, 367, 371, 5, 136,
	69, 2, 368, 371, 5, 104, 53, 2, 369, 371, 5, 98, 50, 2, 370, 365, 3, 2,
	2, 2, 370, 366, 3, 2, 2, 2, 370, 367, 3, 2, 2, 2, 370, 368, 3, 2, 2, 2,
	370, 369, 3, 2, 2, 2, 371, 69, 3, 2, 2, 2, 372, 376, 5, 104, 53, 2, 373,
	376, 5, 138, 70, 2, 374, 376, 5, 98, 50, 2, 375, 372, 3, 2, 2, 2, 375,
	373, 3, 2, 2, 2, 375, 374, 3, 2, 2, 2, 376, 71, 3, 2, 2, 2, 377, 381, 5,
	90, 46, 2, 378, 381, 5, 138, 70, 2, 379, 381, 5, 136, 69, 2, 380, 377,
	3, 2, 2, 2, 380, 378, 3, 2, 2, 2, 380, 379, 3, 2, 2, 2, 381, 382, 3, 2,
	2, 2, 382, 386, 7, 32, 2, 2, 383, 387, 5, 90, 46, 2, 384, 387, 5, 138,
	70, 2, 385, 387, 5, 136, 69, 2, 386, 383, 3, 2, 2, 2, 386, 384, 3, 2, 2,
	2, 386, 385, 3, 2, 2, 2, 387, 73, 3, 2, 2, 2, 388, 400, 7, 11, 2, 2, 389,
	394, 5, 94, 48, 2, 390, 391, 7, 10, 2, 2, 391, 393, 5, 94, 48, 2, 392,
	390, 3, 2, 2, 2, 393, 396, 3, 2, 2, 2, 394, 392, 3, 2, 2, 2, 394, 395,
	3, 2, 2, 2, 395, 398, 3, 2, 2, 2, 396, 394, 3, 2, 2, 2, 397, 399, 7, 10,
	2, 2, 398, 397, 3, 2, 2, 2, 398, 399, 3, 2, 2, 2, 399, 401, 3, 2, 2, 2,
	400, 389, 3, 2, 2, 2, 400, 401, 3, 2, 2, 2, 401, 402, 3, 2, 2, 2, 402,
	403, 7, 12, 2, 2, 403, 75, 3, 2, 2, 2, 404, 416, 7, 15, 2, 2, 405, 410,
	5, 78, 40, 2, 406, 407, 7, 10, 2, 2, 407, 409, 5, 78, 40, 2, 408, 406,
	3, 2, 2, 2, 409, 412, 3, 2, 2, 2, 410, 408, 3, 2, 2, 2, 410, 411, 3, 2,
	2, 2, 411, 414, 3, 2, 2, 2, 412, 410, 3, 2, 2, 2, 413, 415, 7, 10, 2, 2,
	414, 413, 3, 2, 2, 2, 414, 415, 3, 2, 2, 2, 415, 417, 3, 2, 2, 2, 416,
	405, 3, 2, 2, 2, 416, 417, 3, 2, 2, 2, 417, 418, 3, 2, 2, 2, 418, 419,
	7, 16, 2, 2, 419, 77, 3, 2, 2, 2, 420, 421, 5, 82, 42, 2, 421, 422, 7,
	7, 2, 2, 422, 423, 5, 94, 48, 2, 423, 430, 3, 2, 2, 2, 424, 425, 5, 80,
	41, 2, 425, 426, 7, 7, 2, 2, 426, 427, 5, 94, 48, 2, 427, 430, 3, 2, 2,
	2, 428, 430, 5, 138, 70, 2, 429, 420, 3, 2, 2, 2, 429, 424, 3, 2, 2, 2,
	429, 428, 3, 2, 2, 2, 430, 79, 3, 2, 2, 2, 431, 432, 7, 11, 2, 2, 432,
	433, 5, 94, 48, 2, 433, 434, 7, 12, 2, 2, 434, 81, 3, 2, 2, 2, 435, 439,
	7, 66, 2, 2, 436, 439, 5, 86, 44, 2, 437, 439, 5, 136, 69, 2, 438, 435,
	3, 2, 2, 2, 438, 436, 3, 2, 2, 2, 438, 437, 3, 2, 2, 2, 439, 83, 3, 2,
	2, 2, 440, 441, 7, 50, 2, 2, 441, 85, 3, 2, 2, 2, 442, 443, 7, 67, 2, 2,
	443, 87, 3, 2, 2, 2, 444, 445, 7, 69, 2, 2, 445, 89, 3, 2, 2, 2, 446, 447,
	7, 68, 2, 2, 447, 91, 3, 2, 2, 2, 448, 449, 9, 2, 2, 2, 449, 93, 3, 2,
	2, 2, 450, 451, 8, 48, 1, 2, 451, 452, 5, 134, 68, 2, 452, 453, 5, 94,
	48, 27, 453, 469, 3, 2, 2, 2, 454, 469, 5, 72, 37, 2, 455, 469, 5, 86,
	44, 2, 456, 469, 5, 88, 45, 2, 457, 469, 5, 90, 46, 2, 458, 469, 5, 84,
	43, 2, 459, 469, 5, 74, 38, 2, 460, 469, 5, 76, 39, 2, 461, 469, 5, 98,
	50, 2, 462, 469, 5, 104, 53, 2, 463, 469, 5, 136, 69, 2, 464, 469, 5, 138,
	70, 2, 465, 469, 5, 92, 47, 2, 466, 469, 5, 96, 49, 2, 467, 469, 5, 22,
	12, 2, 468, 450, 3, 2, 2, 2, 468, 454, 3, 2, 2, 2, 468, 455, 3, 2, 2, 2,
	468, 456, 3, 2, 2, 2, 468, 457, 3, 2, 2, 2, 468, 458, 3, 2, 2, 2, 468,
	459, 3, 2, 2, 2, 468, 460, 3, 2, 2, 2, 468, 461, 3, 2, 2, 2, 468, 462,
	3, 2, 2, 2, 468, 463, 3, 2, 2, 2, 468, 464, 3, 2, 2, 2, 468, 465, 3, 2,
	2, 2, 468, 466, 3, 2, 2, 2, 468, 467, 3, 2, 2, 2, 469, 519, 3, 2, 2, 2,
	470, 471, 12, 26, 2, 2, 471, 472, 5, 130, 66, 2, 472, 473, 5, 94, 48, 27,
	473, 518, 3, 2, 2, 2, 474, 475, 12, 25, 2, 2, 475, 476, 5, 132, 67, 2,
	476, 477, 5, 94, 48, 26, 477, 518, 3, 2, 2, 2, 478, 479, 12, 24, 2, 2,
	479, 482, 5, 116, 59, 2, 480, 483, 5, 118, 60, 2, 481, 483, 5, 122, 62,
	2, 482, 480, 3, 2, 2, 2, 482, 481, 3, 2, 2, 2, 483, 484, 3, 2, 2, 2, 484,
	485, 5, 94, 48, 25, 485, 518, 3, 2, 2, 2, 486, 487, 12, 23, 2, 2, 487,
	488, 5, 118, 60, 2, 488, 489, 5, 94, 48, 24, 489, 518, 3, 2, 2, 2, 490,
	491, 12, 22, 2, 2, 491, 492, 5, 120, 61, 2, 492, 493, 5, 94, 48, 23, 493,
	518, 3, 2, 2, 2, 494, 495, 12, 21, 2, 2, 495, 496, 5, 122, 62, 2, 496,
	497, 5, 94, 48, 22, 497, 518, 3, 2, 2, 2, 498, 499, 12, 20, 2, 2, 499,
	500, 5, 124, 63, 2, 500, 501, 5, 94, 48, 21, 501, 518, 3, 2, 2, 2, 502,
	503, 12, 19, 2, 2, 503, 504, 5, 126, 64, 2, 504, 505, 5, 94, 48, 20, 505,
	518, 3, 2, 2, 2, 506, 507, 12, 18, 2, 2, 507, 508, 5, 128, 65, 2, 508,
	509, 5, 94, 48, 19, 509, 518, 3, 2, 2, 2, 510, 511, 12, 17, 2, 2, 511,
	513, 7, 34, 2, 2, 512, 514, 5, 94, 48, 2, 513, 512, 3, 2, 2, 2, 513, 514,
	3, 2, 2, 2, 514, 515, 3, 2, 2, 2, 515, 516, 7, 7, 2, 2, 516, 518, 5, 94,
	48, 18, 517, 470, 3, 2, 2, 2, 517, 474, 3, 2, 2, 2, 517, 478, 3, 2, 2,
	2, 517, 486, 3, 2, 2, 2, 517, 490, 3, 2, 2, 2, 517, 494, 3, 2, 2, 2, 517,
	498, 3, 2, 2, 2, 517, 502, 3, 2, 2, 2, 517, 506, 3, 2, 2, 2, 517, 510,
	3, 2, 2, 2, 518, 521, 3, 2, 2, 2, 519, 517, 3, 2, 2, 2, 519, 520, 3, 2,
	2, 2, 520, 95, 3, 2, 2, 2, 521, 519, 3, 2, 2, 2, 522, 523, 7, 13, 2, 2,
	523, 524, 5, 94, 48, 2, 524, 526, 7, 14, 2, 2, 525, 527, 5, 108, 55, 2,
	526, 525, 3, 2, 2, 2, 526, 527, 3, 2, 2, 2, 527, 97, 3, 2, 2, 2, 528, 530,
	5, 100, 51, 2, 529, 531, 5, 106, 54, 2, 530, 529, 3, 2, 2, 2, 531, 532,
	3, 2, 2, 2, 532, 530, 3, 2, 2, 2, 532, 533, 3, 2, 2, 2, 533, 99, 3, 2,
	2, 2, 534, 540, 5, 138, 70, 2, 535, 540, 5, 136, 69, 2, 536, 540, 5, 74,
	38, 2, 537, 540, 5, 76, 39, 2, 538, 540, 5, 102, 52, 2, 539, 534, 3, 2,
	2, 2, 539, 535, 3, 2, 2, 2, 539, 536, 3, 2, 2, 2, 539, 537, 3, 2, 2, 2,
	539, 538, 3, 2, 2, 2, 540, 101, 3, 2, 2, 2, 541, 542, 5, 112, 57, 2, 542,
	543, 5, 110, 56, 2, 543, 544, 5, 114, 58, 2, 544, 103, 3, 2, 2, 2, 545,
	547, 5, 102, 52, 2, 546, 548, 5, 108, 55, 2, 547, 546, 3, 2, 2, 2, 547,
	548, 3, 2, 2, 2, 548, 105, 3, 2, 2, 2, 549, 551, 7, 34, 2, 2, 550, 549,
	3, 2, 2, 2, 550, 551, 3, 2, 2, 2, 551, 552, 3, 2, 2, 2, 552, 553, 7, 9,
	2, 2, 553, 560, 5, 82, 42, 2, 554, 555, 7, 34, 2, 2, 555, 557, 7, 9, 2,
	2, 556, 554, 3, 2, 2, 2, 556, 557, 3, 2, 2, 2, 557, 558, 3, 2, 2, 2, 558,
	560, 5, 80, 41, 2, 559, 550, 3, 2, 2, 2, 559, 556, 3, 2, 2, 2, 560, 107,
	3, 2, 2, 2, 561, 562, 7, 34, 2, 2, 562, 109, 3, 2, 2, 2, 563, 564, 9, 3,
	2, 2, 564, 111, 3, 2, 2, 2, 565, 567, 7, 70, 2, 2, 566, 565, 3, 2, 2, 2,
	567, 570, 3, 2, 2, 2, 568, 566, 3, 2, 2, 2, 568, 569, 3, 2, 2, 2, 569,
	113, 3, 2, 2, 2, 570, 568, 3, 2, 2, 2, 571, 580, 7, 13, 2, 2, 572, 577,
	5, 94, 48, 2, 573, 574, 7, 10, 2, 2, 574, 576, 5, 94, 48, 2, 575, 573,
	3, 2, 2, 2, 576, 579, 3, 2, 2, 2, 577, 575, 3, 2, 2, 2, 577, 578, 3, 2,
	2, 2, 578, 581, 3, 2, 2, 2, 579, 577, 3, 2, 2, 2, 580, 572, 3, 2, 2, 2,
	580, 581, 3, 2, 2, 2, 581, 582, 3, 2, 2, 2, 582, 583, 7, 14, 2, 2, 583,
	115, 3, 2, 2, 2, 584, 585, 9, 4, 2, 2, 585, 117, 3, 2, 2, 2, 586, 590,
	7, 62, 2, 2, 587, 588, 7, 61, 2, 2, 588, 590, 7, 62, 2, 2, 589, 586, 3,
	2, 2, 2, 589, 587, 3, 2, 2, 2, 590, 119, 3, 2, 2, 2, 591, 595, 7, 60, 2,
	2, 592, 593, 7, 61, 2, 2, 593, 595, 7, 60, 2, 2, 594, 591, 3, 2, 2, 2,
	594, 592, 3, 2, 2, 2, 595, 121, 3, 2, 2, 2, 596, 597, 9, 5, 2, 2, 597,
	123, 3, 2, 2, 2, 598, 599, 9, 6, 2, 2, 599, 125, 3, 2, 2, 2, 600, 601,
	7, 30, 2, 2, 601, 127, 3, 2, 2, 2, 602, 603, 7, 31, 2, 2, 603, 129, 3,
	2, 2, 2, 604, 605, 9, 7, 2, 2, 605, 131, 3, 2, 2, 2, 606, 607, 9, 8, 2,
	2, 607, 133, 3, 2, 2, 2, 608, 609, 9, 9, 2, 2, 609, 135, 3, 2, 2, 2, 610,
	611, 7, 65, 2, 2, 611, 612, 7, 66, 2, 2, 612, 137, 3, 2, 2, 2, 613, 614,
	7, 66, 2, 2, 614, 139, 3, 2, 2, 2, 61, 143, 162, 170, 174, 183, 191, 195,
	201, 208, 216, 223, 228, 237, 243, 247, 251, 255, 264, 268, 276, 281, 301,
	312, 321, 334, 336, 347, 355, 358, 363, 370, 375, 380, 386, 394, 398, 400,
	410, 414, 416, 429, 438, 468, 482, 513, 517, 519, 526, 532, 539, 547, 550,
	556, 559, 568, 577, 580, 589, 594,
}
var literalNames = []string{
	"", "", "", "", "", "':'", "';'", "'.'", "','", "'['", "']'", "'('", "')'",
	"'{'", "'}'", "'>'", "'<'", "'=='", "'>='", "'<='", "'!='", "'*'", "'/'",
	"'%'", "'+'", "'-'", "'--'", "'++'", "", "", "", "'='", "'?'", "'!~'",
	"'=~'", "'FOR'", "'RETURN'", "'WAITFOR'", "'OPTIONS'", "'DISTINCT'", "'FILTER'",
	"'SORT'", "'LIMIT'", "'LET'", "'COLLECT'", "", "'NONE'", "'NULL'", "",
	"'USE'", "'INTO'", "'KEEP'", "'WITH'", "'COUNT'", "'ALL'", "'ANY'", "'AGGREGATE'",
	"'EVENT'", "'LIKE'", "", "'IN'", "'DO'", "'WHILE'", "'@'",
}
var symbolicNames = []string{
	"", "MultiLineComment", "SingleLineComment", "WhiteSpaces", "LineTerminator",
	"Colon", "SemiColon", "Dot", "Comma", "OpenBracket", "CloseBracket", "OpenParen",
	"CloseParen", "OpenBrace", "CloseBrace", "Gt", "Lt", "Eq", "Gte", "Lte",
	"Neq", "Multi", "Div", "Mod", "Plus", "Minus", "MinusMinus", "PlusPlus",
	"And", "Or", "Range", "Assign", "QuestionMark", "RegexNotMatch", "RegexMatch",
	"For", "Return", "Waitfor", "Options", "Distinct", "Filter", "Sort", "Limit",
	"Let", "Collect", "SortDirection", "None", "Null", "BooleanLiteral", "Use",
	"Into", "Keep", "With", "Count", "All", "Any", "Aggregate", "Event", "Like",
	"Not", "In", "Do", "While", "Param", "Identifier", "StringLiteral", "IntegerLiteral",
	"FloatLiteral", "NamespaceSegment", "UnknownIdentifier",
}

var ruleNames = []string{
	"program", "head", "useExpression", "use", "namespaceIdentifier", "body",
	"bodyStatement", "bodyExpression", "variableDeclaration", "returnExpression",
	"inlineHighLevelExpression", "highLevelExpression", "forExpression", "forExpressionSource",
	"forExpressionClause", "forExpressionStatement", "forExpressionBody", "forExpressionReturn",
	"filterClause", "limitClause", "limitClauseValue", "sortClause", "sortClauseExpression",
	"collectClause", "collectSelector", "collectGrouping", "collectAggregator",
	"collectAggregateSelector", "collectGroupVariable", "collectCounter", "optionsClause",
	"waitForExpression", "waitForTimeout", "waitForEventName", "waitForEventSource",
	"rangeOperator", "arrayLiteral", "objectLiteral", "propertyAssignment",
	"computedPropertyName", "propertyName", "booleanLiteral", "stringLiteral",
	"floatLiteral", "integerLiteral", "noneLiteral", "expression", "expressionGroup",
	"memberExpression", "memberExpressionSource", "functionCall", "functionCallExpression",
	"memberExpressionPath", "errorOperator", "functionIdentifier", "namespace",
	"arguments", "arrayOperator", "inOperator", "likeOperator", "equalityOperator",
	"regexpOperator", "logicalAndOperator", "logicalOrOperator", "multiplicativeOperator",
	"additiveOperator", "unaryOperator", "param", "variable",
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
	FqlParserDistinct          = 39
	FqlParserFilter            = 40
	FqlParserSort              = 41
	FqlParserLimit             = 42
	FqlParserLet               = 43
	FqlParserCollect           = 44
	FqlParserSortDirection     = 45
	FqlParserNone              = 46
	FqlParserNull              = 47
	FqlParserBooleanLiteral    = 48
	FqlParserUse               = 49
	FqlParserInto              = 50
	FqlParserKeep              = 51
	FqlParserWith              = 52
	FqlParserCount             = 53
	FqlParserAll               = 54
	FqlParserAny               = 55
	FqlParserAggregate         = 56
	FqlParserEvent             = 57
	FqlParserLike              = 58
	FqlParserNot               = 59
	FqlParserIn                = 60
	FqlParserDo                = 61
	FqlParserWhile             = 62
	FqlParserParam             = 63
	FqlParserIdentifier        = 64
	FqlParserStringLiteral     = 65
	FqlParserIntegerLiteral    = 66
	FqlParserFloatLiteral      = 67
	FqlParserNamespaceSegment  = 68
	FqlParserUnknownIdentifier = 69
)

// FqlParser rules.
const (
	FqlParserRULE_program                   = 0
	FqlParserRULE_head                      = 1
	FqlParserRULE_useExpression             = 2
	FqlParserRULE_use                       = 3
	FqlParserRULE_namespaceIdentifier       = 4
	FqlParserRULE_body                      = 5
	FqlParserRULE_bodyStatement             = 6
	FqlParserRULE_bodyExpression            = 7
	FqlParserRULE_variableDeclaration       = 8
	FqlParserRULE_returnExpression          = 9
	FqlParserRULE_inlineHighLevelExpression = 10
	FqlParserRULE_highLevelExpression       = 11
	FqlParserRULE_forExpression             = 12
	FqlParserRULE_forExpressionSource       = 13
	FqlParserRULE_forExpressionClause       = 14
	FqlParserRULE_forExpressionStatement    = 15
	FqlParserRULE_forExpressionBody         = 16
	FqlParserRULE_forExpressionReturn       = 17
	FqlParserRULE_filterClause              = 18
	FqlParserRULE_limitClause               = 19
	FqlParserRULE_limitClauseValue          = 20
	FqlParserRULE_sortClause                = 21
	FqlParserRULE_sortClauseExpression      = 22
	FqlParserRULE_collectClause             = 23
	FqlParserRULE_collectSelector           = 24
	FqlParserRULE_collectGrouping           = 25
	FqlParserRULE_collectAggregator         = 26
	FqlParserRULE_collectAggregateSelector  = 27
	FqlParserRULE_collectGroupVariable      = 28
	FqlParserRULE_collectCounter            = 29
	FqlParserRULE_optionsClause             = 30
	FqlParserRULE_waitForExpression         = 31
	FqlParserRULE_waitForTimeout            = 32
	FqlParserRULE_waitForEventName          = 33
	FqlParserRULE_waitForEventSource        = 34
	FqlParserRULE_rangeOperator             = 35
	FqlParserRULE_arrayLiteral              = 36
	FqlParserRULE_objectLiteral             = 37
	FqlParserRULE_propertyAssignment        = 38
	FqlParserRULE_computedPropertyName      = 39
	FqlParserRULE_propertyName              = 40
	FqlParserRULE_booleanLiteral            = 41
	FqlParserRULE_stringLiteral             = 42
	FqlParserRULE_floatLiteral              = 43
	FqlParserRULE_integerLiteral            = 44
	FqlParserRULE_noneLiteral               = 45
	FqlParserRULE_expression                = 46
	FqlParserRULE_expressionGroup           = 47
	FqlParserRULE_memberExpression          = 48
	FqlParserRULE_memberExpressionSource    = 49
	FqlParserRULE_functionCall              = 50
	FqlParserRULE_functionCallExpression    = 51
	FqlParserRULE_memberExpressionPath      = 52
	FqlParserRULE_errorOperator             = 53
	FqlParserRULE_functionIdentifier        = 54
	FqlParserRULE_namespace                 = 55
	FqlParserRULE_arguments                 = 56
	FqlParserRULE_arrayOperator             = 57
	FqlParserRULE_inOperator                = 58
	FqlParserRULE_likeOperator              = 59
	FqlParserRULE_equalityOperator          = 60
	FqlParserRULE_regexpOperator            = 61
	FqlParserRULE_logicalAndOperator        = 62
	FqlParserRULE_logicalOrOperator         = 63
	FqlParserRULE_multiplicativeOperator    = 64
	FqlParserRULE_additiveOperator          = 65
	FqlParserRULE_unaryOperator             = 66
	FqlParserRULE_param                     = 67
	FqlParserRULE_variable                  = 68
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

func (s *ProgramContext) EOF() antlr.TerminalNode {
	return s.GetToken(FqlParserEOF, 0)
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
	{
		p.SetState(145)
		p.Match(FqlParserEOF)
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
		p.SetState(147)
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
		p.SetState(149)
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
		p.SetState(151)
		p.Match(FqlParserUse)
	}
	{
		p.SetState(152)
		p.NamespaceIdentifier()
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
	p.EnterRule(localctx, 8, FqlParserRULE_namespaceIdentifier)

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
		p.Namespace()
	}
	{
		p.SetState(155)
		p.Match(FqlParserIdentifier)
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
	p.EnterRule(localctx, 10, FqlParserRULE_body)

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
	p.SetState(160)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 1, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			{
				p.SetState(157)
				p.BodyStatement()
			}

		}
		p.SetState(162)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 1, p.GetParserRuleContext())
	}
	{
		p.SetState(163)
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
	localctx = NewBodyStatementContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 12, FqlParserRULE_bodyStatement)

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
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 2, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(165)
			p.VariableDeclaration()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(166)
			p.FunctionCallExpression()
		}

	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(167)
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
	localctx = NewBodyExpressionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 14, FqlParserRULE_bodyExpression)

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

	p.SetState(172)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case FqlParserReturn:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(170)
			p.ReturnExpression()
		}

	case FqlParserFor:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(171)
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
	p.EnterRule(localctx, 16, FqlParserRULE_variableDeclaration)

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
		p.SetState(174)
		p.Match(FqlParserLet)
	}
	{
		p.SetState(175)
		p.Match(FqlParserIdentifier)
	}
	{
		p.SetState(176)
		p.Match(FqlParserAssign)
	}
	{
		p.SetState(177)
		p.expression(0)
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
	localctx = NewReturnExpressionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 18, FqlParserRULE_returnExpression)

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
		p.SetState(179)
		p.Match(FqlParserReturn)
	}
	p.SetState(181)
	p.GetErrorHandler().Sync(p)

	if p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 4, p.GetParserRuleContext()) == 1 {
		{
			p.SetState(180)
			p.Match(FqlParserDistinct)
		}

	}
	{
		p.SetState(183)
		p.expression(0)
	}

	return localctx
}

// IInlineHighLevelExpressionContext is an interface to support dynamic dispatch.
type IInlineHighLevelExpressionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsInlineHighLevelExpressionContext differentiates from other interfaces.
	IsInlineHighLevelExpressionContext()
}

type InlineHighLevelExpressionContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyInlineHighLevelExpressionContext() *InlineHighLevelExpressionContext {
	var p = new(InlineHighLevelExpressionContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = FqlParserRULE_inlineHighLevelExpression
	return p
}

func (*InlineHighLevelExpressionContext) IsInlineHighLevelExpressionContext() {}

func NewInlineHighLevelExpressionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *InlineHighLevelExpressionContext {
	var p = new(InlineHighLevelExpressionContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_inlineHighLevelExpression

	return p
}

func (s *InlineHighLevelExpressionContext) GetParser() antlr.Parser { return s.parser }

func (s *InlineHighLevelExpressionContext) OpenParen() antlr.TerminalNode {
	return s.GetToken(FqlParserOpenParen, 0)
}

func (s *InlineHighLevelExpressionContext) HighLevelExpression() IHighLevelExpressionContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IHighLevelExpressionContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IHighLevelExpressionContext)
}

func (s *InlineHighLevelExpressionContext) CloseParen() antlr.TerminalNode {
	return s.GetToken(FqlParserCloseParen, 0)
}

func (s *InlineHighLevelExpressionContext) ErrorOperator() IErrorOperatorContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IErrorOperatorContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IErrorOperatorContext)
}

func (s *InlineHighLevelExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *InlineHighLevelExpressionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *InlineHighLevelExpressionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterInlineHighLevelExpression(s)
	}
}

func (s *InlineHighLevelExpressionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitInlineHighLevelExpression(s)
	}
}

func (s *InlineHighLevelExpressionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitInlineHighLevelExpression(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) InlineHighLevelExpression() (localctx IInlineHighLevelExpressionContext) {
	localctx = NewInlineHighLevelExpressionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 20, FqlParserRULE_inlineHighLevelExpression)

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
		p.SetState(185)
		p.Match(FqlParserOpenParen)
	}
	{
		p.SetState(186)
		p.HighLevelExpression()
	}
	{
		p.SetState(187)
		p.Match(FqlParserCloseParen)
	}
	p.SetState(189)
	p.GetErrorHandler().Sync(p)

	if p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 5, p.GetParserRuleContext()) == 1 {
		{
			p.SetState(188)
			p.ErrorOperator()
		}

	}

	return localctx
}

// IHighLevelExpressionContext is an interface to support dynamic dispatch.
type IHighLevelExpressionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsHighLevelExpressionContext differentiates from other interfaces.
	IsHighLevelExpressionContext()
}

type HighLevelExpressionContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyHighLevelExpressionContext() *HighLevelExpressionContext {
	var p = new(HighLevelExpressionContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = FqlParserRULE_highLevelExpression
	return p
}

func (*HighLevelExpressionContext) IsHighLevelExpressionContext() {}

func NewHighLevelExpressionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *HighLevelExpressionContext {
	var p = new(HighLevelExpressionContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_highLevelExpression

	return p
}

func (s *HighLevelExpressionContext) GetParser() antlr.Parser { return s.parser }

func (s *HighLevelExpressionContext) ForExpression() IForExpressionContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IForExpressionContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IForExpressionContext)
}

func (s *HighLevelExpressionContext) WaitForExpression() IWaitForExpressionContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IWaitForExpressionContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IWaitForExpressionContext)
}

func (s *HighLevelExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *HighLevelExpressionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *HighLevelExpressionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterHighLevelExpression(s)
	}
}

func (s *HighLevelExpressionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitHighLevelExpression(s)
	}
}

func (s *HighLevelExpressionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitHighLevelExpression(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) HighLevelExpression() (localctx IHighLevelExpressionContext) {
	localctx = NewHighLevelExpressionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 22, FqlParserRULE_highLevelExpression)

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

	p.SetState(193)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case FqlParserFor:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(191)
			p.ForExpression()
		}

	case FqlParserWaitfor:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(192)
			p.WaitForExpression()
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
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

func (s *ForExpressionContext) AllIdentifier() []antlr.TerminalNode {
	return s.GetTokens(FqlParserIdentifier)
}

func (s *ForExpressionContext) Identifier(i int) antlr.TerminalNode {
	return s.GetToken(FqlParserIdentifier, i)
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
	p.EnterRule(localctx, 24, FqlParserRULE_forExpression)
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

	p.SetState(226)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 11, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(195)
			p.Match(FqlParserFor)
		}
		{
			p.SetState(196)
			p.Match(FqlParserIdentifier)
		}
		p.SetState(199)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		if _la == FqlParserComma {
			{
				p.SetState(197)
				p.Match(FqlParserComma)
			}
			{
				p.SetState(198)
				p.Match(FqlParserIdentifier)
			}

		}
		{
			p.SetState(201)
			p.Match(FqlParserIn)
		}
		{
			p.SetState(202)
			p.ForExpressionSource()
		}
		p.SetState(206)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 8, p.GetParserRuleContext())

		for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
			if _alt == 1 {
				{
					p.SetState(203)
					p.ForExpressionBody()
				}

			}
			p.SetState(208)
			p.GetErrorHandler().Sync(p)
			_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 8, p.GetParserRuleContext())
		}
		{
			p.SetState(209)
			p.ForExpressionReturn()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(211)
			p.Match(FqlParserFor)
		}
		{
			p.SetState(212)
			p.Match(FqlParserIdentifier)
		}
		p.SetState(214)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		if _la == FqlParserDo {
			{
				p.SetState(213)
				p.Match(FqlParserDo)
			}

		}
		{
			p.SetState(216)
			p.Match(FqlParserWhile)
		}
		{
			p.SetState(217)
			p.expression(0)
		}
		p.SetState(221)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 10, p.GetParserRuleContext())

		for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
			if _alt == 1 {
				{
					p.SetState(218)
					p.ForExpressionBody()
				}

			}
			p.SetState(223)
			p.GetErrorHandler().Sync(p)
			_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 10, p.GetParserRuleContext())
		}
		{
			p.SetState(224)
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
	localctx = NewForExpressionSourceContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 26, FqlParserRULE_forExpressionSource)

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

	p.SetState(235)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 12, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(228)
			p.FunctionCallExpression()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(229)
			p.ArrayLiteral()
		}

	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(230)
			p.ObjectLiteral()
		}

	case 4:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(231)
			p.Variable()
		}

	case 5:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(232)
			p.MemberExpression()
		}

	case 6:
		p.EnterOuterAlt(localctx, 6)
		{
			p.SetState(233)
			p.RangeOperator()
		}

	case 7:
		p.EnterOuterAlt(localctx, 7)
		{
			p.SetState(234)
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
	p.EnterRule(localctx, 28, FqlParserRULE_forExpressionClause)

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

	switch p.GetTokenStream().LA(1) {
	case FqlParserLimit:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(237)
			p.LimitClause()
		}

	case FqlParserSort:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(238)
			p.SortClause()
		}

	case FqlParserFilter:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(239)
			p.FilterClause()
		}

	case FqlParserCollect:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(240)
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
	p.EnterRule(localctx, 30, FqlParserRULE_forExpressionStatement)

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
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 14, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(243)
			p.VariableDeclaration()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(244)
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
	p.EnterRule(localctx, 32, FqlParserRULE_forExpressionBody)

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
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 15, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(247)
			p.ForExpressionStatement()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(248)
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
	p.EnterRule(localctx, 34, FqlParserRULE_forExpressionReturn)

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

	p.SetState(253)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case FqlParserReturn:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(251)
			p.ReturnExpression()
		}

	case FqlParserFor:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(252)
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
	p.EnterRule(localctx, 36, FqlParserRULE_filterClause)

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
		p.Match(FqlParserFilter)
	}
	{
		p.SetState(256)
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
	p.EnterRule(localctx, 38, FqlParserRULE_limitClause)
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
		p.SetState(258)
		p.Match(FqlParserLimit)
	}
	{
		p.SetState(259)
		p.LimitClauseValue()
	}
	p.SetState(262)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == FqlParserComma {
		{
			p.SetState(260)
			p.Match(FqlParserComma)
		}
		{
			p.SetState(261)
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
	p.EnterRule(localctx, 40, FqlParserRULE_limitClauseValue)

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

	p.SetState(266)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case FqlParserIntegerLiteral:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(264)
			p.Match(FqlParserIntegerLiteral)
		}

	case FqlParserParam:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(265)
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
	p.EnterRule(localctx, 42, FqlParserRULE_sortClause)
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
		p.SetState(268)
		p.Match(FqlParserSort)
	}
	{
		p.SetState(269)
		p.SortClauseExpression()
	}
	p.SetState(274)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == FqlParserComma {
		{
			p.SetState(270)
			p.Match(FqlParserComma)
		}
		{
			p.SetState(271)
			p.SortClauseExpression()
		}

		p.SetState(276)
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
	p.EnterRule(localctx, 44, FqlParserRULE_sortClauseExpression)

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
		p.SetState(277)
		p.expression(0)
	}
	p.SetState(279)
	p.GetErrorHandler().Sync(p)

	if p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 20, p.GetParserRuleContext()) == 1 {
		{
			p.SetState(278)
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
	p.EnterRule(localctx, 46, FqlParserRULE_collectClause)

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

	p.SetState(299)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 21, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(281)
			p.Match(FqlParserCollect)
		}
		{
			p.SetState(282)
			p.CollectCounter()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(283)
			p.Match(FqlParserCollect)
		}
		{
			p.SetState(284)
			p.CollectAggregator()
		}

	case 3:
		p.EnterOuterAlt(localctx, 3)
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
			p.CollectAggregator()
		}

	case 4:
		p.EnterOuterAlt(localctx, 4)
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
			p.CollectGroupVariable()
		}

	case 5:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(293)
			p.Match(FqlParserCollect)
		}
		{
			p.SetState(294)
			p.CollectGrouping()
		}
		{
			p.SetState(295)
			p.CollectCounter()
		}

	case 6:
		p.EnterOuterAlt(localctx, 6)
		{
			p.SetState(297)
			p.Match(FqlParserCollect)
		}
		{
			p.SetState(298)
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
	p.EnterRule(localctx, 48, FqlParserRULE_collectSelector)

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
		p.Match(FqlParserIdentifier)
	}
	{
		p.SetState(302)
		p.Match(FqlParserAssign)
	}
	{
		p.SetState(303)
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
	p.EnterRule(localctx, 50, FqlParserRULE_collectGrouping)
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
		p.SetState(305)
		p.CollectSelector()
	}
	p.SetState(310)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == FqlParserComma {
		{
			p.SetState(306)
			p.Match(FqlParserComma)
		}
		{
			p.SetState(307)
			p.CollectSelector()
		}

		p.SetState(312)
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
	p.EnterRule(localctx, 52, FqlParserRULE_collectAggregator)
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
		p.SetState(313)
		p.Match(FqlParserAggregate)
	}
	{
		p.SetState(314)
		p.CollectAggregateSelector()
	}
	p.SetState(319)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == FqlParserComma {
		{
			p.SetState(315)
			p.Match(FqlParserComma)
		}
		{
			p.SetState(316)
			p.CollectAggregateSelector()
		}

		p.SetState(321)
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
	p.EnterRule(localctx, 54, FqlParserRULE_collectAggregateSelector)

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
		p.SetState(322)
		p.Match(FqlParserIdentifier)
	}
	{
		p.SetState(323)
		p.Match(FqlParserAssign)
	}
	{
		p.SetState(324)
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
	p.EnterRule(localctx, 56, FqlParserRULE_collectGroupVariable)

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

	p.SetState(334)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 25, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(326)
			p.Match(FqlParserInto)
		}
		{
			p.SetState(327)
			p.CollectSelector()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(328)
			p.Match(FqlParserInto)
		}
		{
			p.SetState(329)
			p.Match(FqlParserIdentifier)
		}
		p.SetState(332)
		p.GetErrorHandler().Sync(p)

		if p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 24, p.GetParserRuleContext()) == 1 {
			{
				p.SetState(330)
				p.Match(FqlParserKeep)
			}
			{
				p.SetState(331)
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
	p.EnterRule(localctx, 58, FqlParserRULE_collectCounter)

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
		p.SetState(336)
		p.Match(FqlParserWith)
	}
	{
		p.SetState(337)
		p.Match(FqlParserCount)
	}
	{
		p.SetState(338)
		p.Match(FqlParserInto)
	}
	{
		p.SetState(339)
		p.Match(FqlParserIdentifier)
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

func (s *OptionsClauseContext) Variable() IVariableContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IVariableContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IVariableContext)
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

	p.SetState(345)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 26, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(341)
			p.Match(FqlParserOptions)
		}
		{
			p.SetState(342)
			p.ObjectLiteral()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(343)
			p.Match(FqlParserOptions)
		}
		{
			p.SetState(344)
			p.Variable()
		}

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

func (s *WaitForExpressionContext) WaitForTimeout() IWaitForTimeoutContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IWaitForTimeoutContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IWaitForTimeoutContext)
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
	localctx = NewWaitForExpressionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 62, FqlParserRULE_waitForExpression)
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
		p.SetState(347)
		p.Match(FqlParserWaitfor)
	}
	{
		p.SetState(348)
		p.Match(FqlParserEvent)
	}
	{
		p.SetState(349)
		p.WaitForEventName()
	}
	{
		p.SetState(350)
		p.Match(FqlParserIn)
	}
	{
		p.SetState(351)
		p.WaitForEventSource()
	}
	p.SetState(353)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == FqlParserOptions {
		{
			p.SetState(352)
			p.OptionsClause()
		}

	}
	p.SetState(356)
	p.GetErrorHandler().Sync(p)

	if p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 28, p.GetParserRuleContext()) == 1 {
		{
			p.SetState(355)
			p.WaitForTimeout()
		}

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

func (s *WaitForTimeoutContext) Param() IParamContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IParamContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IParamContext)
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
	p.EnterRule(localctx, 64, FqlParserRULE_waitForTimeout)

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

	switch p.GetTokenStream().LA(1) {
	case FqlParserIntegerLiteral:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(358)
			p.IntegerLiteral()
		}

	case FqlParserIdentifier:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(359)
			p.Variable()
		}

	case FqlParserParam:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(360)
			p.Param()
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
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
	p.EnterRule(localctx, 66, FqlParserRULE_waitForEventName)

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

	p.SetState(368)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 30, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(363)
			p.StringLiteral()
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
			p.Param()
		}

	case 4:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(366)
			p.FunctionCallExpression()
		}

	case 5:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(367)
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
	localctx = NewWaitForEventSourceContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 68, FqlParserRULE_waitForEventSource)

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

	p.SetState(373)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 31, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(370)
			p.FunctionCallExpression()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(371)
			p.Variable()
		}

	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(372)
			p.MemberExpression()
		}

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
	p.EnterRule(localctx, 70, FqlParserRULE_rangeOperator)

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
	p.SetState(378)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case FqlParserIntegerLiteral:
		{
			p.SetState(375)
			p.IntegerLiteral()
		}

	case FqlParserIdentifier:
		{
			p.SetState(376)
			p.Variable()
		}

	case FqlParserParam:
		{
			p.SetState(377)
			p.Param()
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}
	{
		p.SetState(380)
		p.Match(FqlParserRange)
	}
	p.SetState(384)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case FqlParserIntegerLiteral:
		{
			p.SetState(381)
			p.IntegerLiteral()
		}

	case FqlParserIdentifier:
		{
			p.SetState(382)
			p.Variable()
		}

	case FqlParserParam:
		{
			p.SetState(383)
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

func (s *ArrayLiteralContext) AllExpression() []IExpressionContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IExpressionContext)(nil)).Elem())
	var tst = make([]IExpressionContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IExpressionContext)
		}
	}

	return tst
}

func (s *ArrayLiteralContext) Expression(i int) IExpressionContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExpressionContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *ArrayLiteralContext) AllComma() []antlr.TerminalNode {
	return s.GetTokens(FqlParserComma)
}

func (s *ArrayLiteralContext) Comma(i int) antlr.TerminalNode {
	return s.GetToken(FqlParserComma, i)
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
	p.EnterRule(localctx, 72, FqlParserRULE_arrayLiteral)
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
		p.SetState(386)
		p.Match(FqlParserOpenBracket)
	}
	p.SetState(398)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if (((_la-9)&-(0x1f+1)) == 0 && ((1<<uint((_la-9)))&((1<<(FqlParserOpenBracket-9))|(1<<(FqlParserOpenParen-9))|(1<<(FqlParserOpenBrace-9))|(1<<(FqlParserPlus-9))|(1<<(FqlParserMinus-9))|(1<<(FqlParserAnd-9))|(1<<(FqlParserOr-9))|(1<<(FqlParserFor-9))|(1<<(FqlParserReturn-9))|(1<<(FqlParserWaitfor-9))|(1<<(FqlParserDistinct-9))|(1<<(FqlParserFilter-9)))) != 0) || (((_la-41)&-(0x1f+1)) == 0 && ((1<<uint((_la-41)))&((1<<(FqlParserSort-41))|(1<<(FqlParserLimit-41))|(1<<(FqlParserLet-41))|(1<<(FqlParserCollect-41))|(1<<(FqlParserSortDirection-41))|(1<<(FqlParserNone-41))|(1<<(FqlParserNull-41))|(1<<(FqlParserBooleanLiteral-41))|(1<<(FqlParserUse-41))|(1<<(FqlParserInto-41))|(1<<(FqlParserKeep-41))|(1<<(FqlParserWith-41))|(1<<(FqlParserCount-41))|(1<<(FqlParserAll-41))|(1<<(FqlParserAny-41))|(1<<(FqlParserAggregate-41))|(1<<(FqlParserEvent-41))|(1<<(FqlParserLike-41))|(1<<(FqlParserNot-41))|(1<<(FqlParserIn-41))|(1<<(FqlParserParam-41))|(1<<(FqlParserIdentifier-41))|(1<<(FqlParserStringLiteral-41))|(1<<(FqlParserIntegerLiteral-41))|(1<<(FqlParserFloatLiteral-41))|(1<<(FqlParserNamespaceSegment-41)))) != 0) {
		{
			p.SetState(387)
			p.expression(0)
		}
		p.SetState(392)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 34, p.GetParserRuleContext())

		for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
			if _alt == 1 {
				{
					p.SetState(388)
					p.Match(FqlParserComma)
				}
				{
					p.SetState(389)
					p.expression(0)
				}

			}
			p.SetState(394)
			p.GetErrorHandler().Sync(p)
			_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 34, p.GetParserRuleContext())
		}
		p.SetState(396)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		if _la == FqlParserComma {
			{
				p.SetState(395)
				p.Match(FqlParserComma)
			}

		}

	}
	{
		p.SetState(400)
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
	p.EnterRule(localctx, 74, FqlParserRULE_objectLiteral)
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
		p.SetState(402)
		p.Match(FqlParserOpenBrace)
	}
	p.SetState(414)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == FqlParserOpenBracket || (((_la-63)&-(0x1f+1)) == 0 && ((1<<uint((_la-63)))&((1<<(FqlParserParam-63))|(1<<(FqlParserIdentifier-63))|(1<<(FqlParserStringLiteral-63)))) != 0) {
		{
			p.SetState(403)
			p.PropertyAssignment()
		}
		p.SetState(408)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 37, p.GetParserRuleContext())

		for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
			if _alt == 1 {
				{
					p.SetState(404)
					p.Match(FqlParserComma)
				}
				{
					p.SetState(405)
					p.PropertyAssignment()
				}

			}
			p.SetState(410)
			p.GetErrorHandler().Sync(p)
			_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 37, p.GetParserRuleContext())
		}
		p.SetState(412)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		if _la == FqlParserComma {
			{
				p.SetState(411)
				p.Match(FqlParserComma)
			}

		}

	}
	{
		p.SetState(416)
		p.Match(FqlParserCloseBrace)
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
	localctx = NewPropertyAssignmentContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 76, FqlParserRULE_propertyAssignment)

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

	p.SetState(427)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 40, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(418)
			p.PropertyName()
		}
		{
			p.SetState(419)
			p.Match(FqlParserColon)
		}
		{
			p.SetState(420)
			p.expression(0)
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(422)
			p.ComputedPropertyName()
		}
		{
			p.SetState(423)
			p.Match(FqlParserColon)
		}
		{
			p.SetState(424)
			p.expression(0)
		}

	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(426)
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
	localctx = NewComputedPropertyNameContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 78, FqlParserRULE_computedPropertyName)

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
		p.SetState(429)
		p.Match(FqlParserOpenBracket)
	}
	{
		p.SetState(430)
		p.expression(0)
	}
	{
		p.SetState(431)
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
	p.EnterRule(localctx, 80, FqlParserRULE_propertyName)

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

	p.SetState(436)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case FqlParserIdentifier:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(433)
			p.Match(FqlParserIdentifier)
		}

	case FqlParserStringLiteral:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(434)
			p.StringLiteral()
		}

	case FqlParserParam:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(435)
			p.Param()
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
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
	p.EnterRule(localctx, 82, FqlParserRULE_booleanLiteral)

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
		p.SetState(438)
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
	p.EnterRule(localctx, 84, FqlParserRULE_stringLiteral)

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
		p.SetState(440)
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
	localctx = NewFloatLiteralContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 86, FqlParserRULE_floatLiteral)

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
		p.SetState(442)
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
	localctx = NewIntegerLiteralContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 88, FqlParserRULE_integerLiteral)

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
		p.SetState(444)
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
	localctx = NewNoneLiteralContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 90, FqlParserRULE_noneLiteral)
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
		p.SetState(446)
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

func (s *ExpressionContext) FloatLiteral() IFloatLiteralContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IFloatLiteralContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IFloatLiteralContext)
}

func (s *ExpressionContext) IntegerLiteral() IIntegerLiteralContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IIntegerLiteralContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IIntegerLiteralContext)
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

func (s *ExpressionContext) MemberExpression() IMemberExpressionContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IMemberExpressionContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IMemberExpressionContext)
}

func (s *ExpressionContext) FunctionCallExpression() IFunctionCallExpressionContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IFunctionCallExpressionContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IFunctionCallExpressionContext)
}

func (s *ExpressionContext) Param() IParamContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IParamContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IParamContext)
}

func (s *ExpressionContext) Variable() IVariableContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IVariableContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IVariableContext)
}

func (s *ExpressionContext) NoneLiteral() INoneLiteralContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*INoneLiteralContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(INoneLiteralContext)
}

func (s *ExpressionContext) ExpressionGroup() IExpressionGroupContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExpressionGroupContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IExpressionGroupContext)
}

func (s *ExpressionContext) InlineHighLevelExpression() IInlineHighLevelExpressionContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IInlineHighLevelExpressionContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IInlineHighLevelExpressionContext)
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

func (s *ExpressionContext) LikeOperator() ILikeOperatorContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ILikeOperatorContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ILikeOperatorContext)
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
	_startState := 92
	p.EnterRecursionRule(localctx, 92, FqlParserRULE_expression, _p)
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
	p.SetState(466)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 42, p.GetParserRuleContext()) {
	case 1:
		{
			p.SetState(449)
			p.UnaryOperator()
		}
		{
			p.SetState(450)
			p.expression(25)
		}

	case 2:
		{
			p.SetState(452)
			p.RangeOperator()
		}

	case 3:
		{
			p.SetState(453)
			p.StringLiteral()
		}

	case 4:
		{
			p.SetState(454)
			p.FloatLiteral()
		}

	case 5:
		{
			p.SetState(455)
			p.IntegerLiteral()
		}

	case 6:
		{
			p.SetState(456)
			p.BooleanLiteral()
		}

	case 7:
		{
			p.SetState(457)
			p.ArrayLiteral()
		}

	case 8:
		{
			p.SetState(458)
			p.ObjectLiteral()
		}

	case 9:
		{
			p.SetState(459)
			p.MemberExpression()
		}

	case 10:
		{
			p.SetState(460)
			p.FunctionCallExpression()
		}

	case 11:
		{
			p.SetState(461)
			p.Param()
		}

	case 12:
		{
			p.SetState(462)
			p.Variable()
		}

	case 13:
		{
			p.SetState(463)
			p.NoneLiteral()
		}

	case 14:
		{
			p.SetState(464)
			p.ExpressionGroup()
		}

	case 15:
		{
			p.SetState(465)
			p.InlineHighLevelExpression()
		}

	}
	p.GetParserRuleContext().SetStop(p.GetTokenStream().LT(-1))
	p.SetState(517)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 46, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			if p.GetParseListeners() != nil {
				p.TriggerExitRuleEvent()
			}
			_prevctx = localctx
			p.SetState(515)
			p.GetErrorHandler().Sync(p)
			switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 45, p.GetParserRuleContext()) {
			case 1:
				localctx = NewExpressionContext(p, _parentctx, _parentState)
				p.PushNewRecursionContext(localctx, _startState, FqlParserRULE_expression)
				p.SetState(468)

				if !(p.Precpred(p.GetParserRuleContext(), 24)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 24)", ""))
				}
				{
					p.SetState(469)
					p.MultiplicativeOperator()
				}
				{
					p.SetState(470)
					p.expression(25)
				}

			case 2:
				localctx = NewExpressionContext(p, _parentctx, _parentState)
				p.PushNewRecursionContext(localctx, _startState, FqlParserRULE_expression)
				p.SetState(472)

				if !(p.Precpred(p.GetParserRuleContext(), 23)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 23)", ""))
				}
				{
					p.SetState(473)
					p.AdditiveOperator()
				}
				{
					p.SetState(474)
					p.expression(24)
				}

			case 3:
				localctx = NewExpressionContext(p, _parentctx, _parentState)
				p.PushNewRecursionContext(localctx, _startState, FqlParserRULE_expression)
				p.SetState(476)

				if !(p.Precpred(p.GetParserRuleContext(), 22)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 22)", ""))
				}
				{
					p.SetState(477)
					p.ArrayOperator()
				}
				p.SetState(480)
				p.GetErrorHandler().Sync(p)

				switch p.GetTokenStream().LA(1) {
				case FqlParserNot, FqlParserIn:
					{
						p.SetState(478)
						p.InOperator()
					}

				case FqlParserGt, FqlParserLt, FqlParserEq, FqlParserGte, FqlParserLte, FqlParserNeq:
					{
						p.SetState(479)
						p.EqualityOperator()
					}

				default:
					panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
				}
				{
					p.SetState(482)
					p.expression(23)
				}

			case 4:
				localctx = NewExpressionContext(p, _parentctx, _parentState)
				p.PushNewRecursionContext(localctx, _startState, FqlParserRULE_expression)
				p.SetState(484)

				if !(p.Precpred(p.GetParserRuleContext(), 21)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 21)", ""))
				}
				{
					p.SetState(485)
					p.InOperator()
				}
				{
					p.SetState(486)
					p.expression(22)
				}

			case 5:
				localctx = NewExpressionContext(p, _parentctx, _parentState)
				p.PushNewRecursionContext(localctx, _startState, FqlParserRULE_expression)
				p.SetState(488)

				if !(p.Precpred(p.GetParserRuleContext(), 20)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 20)", ""))
				}
				{
					p.SetState(489)
					p.LikeOperator()
				}
				{
					p.SetState(490)
					p.expression(21)
				}

			case 6:
				localctx = NewExpressionContext(p, _parentctx, _parentState)
				p.PushNewRecursionContext(localctx, _startState, FqlParserRULE_expression)
				p.SetState(492)

				if !(p.Precpred(p.GetParserRuleContext(), 19)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 19)", ""))
				}
				{
					p.SetState(493)
					p.EqualityOperator()
				}
				{
					p.SetState(494)
					p.expression(20)
				}

			case 7:
				localctx = NewExpressionContext(p, _parentctx, _parentState)
				p.PushNewRecursionContext(localctx, _startState, FqlParserRULE_expression)
				p.SetState(496)

				if !(p.Precpred(p.GetParserRuleContext(), 18)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 18)", ""))
				}
				{
					p.SetState(497)
					p.RegexpOperator()
				}
				{
					p.SetState(498)
					p.expression(19)
				}

			case 8:
				localctx = NewExpressionContext(p, _parentctx, _parentState)
				p.PushNewRecursionContext(localctx, _startState, FqlParserRULE_expression)
				p.SetState(500)

				if !(p.Precpred(p.GetParserRuleContext(), 17)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 17)", ""))
				}
				{
					p.SetState(501)
					p.LogicalAndOperator()
				}
				{
					p.SetState(502)
					p.expression(18)
				}

			case 9:
				localctx = NewExpressionContext(p, _parentctx, _parentState)
				p.PushNewRecursionContext(localctx, _startState, FqlParserRULE_expression)
				p.SetState(504)

				if !(p.Precpred(p.GetParserRuleContext(), 16)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 16)", ""))
				}
				{
					p.SetState(505)
					p.LogicalOrOperator()
				}
				{
					p.SetState(506)
					p.expression(17)
				}

			case 10:
				localctx = NewExpressionContext(p, _parentctx, _parentState)
				p.PushNewRecursionContext(localctx, _startState, FqlParserRULE_expression)
				p.SetState(508)

				if !(p.Precpred(p.GetParserRuleContext(), 15)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 15)", ""))
				}
				{
					p.SetState(509)
					p.Match(FqlParserQuestionMark)
				}
				p.SetState(511)
				p.GetErrorHandler().Sync(p)
				_la = p.GetTokenStream().LA(1)

				if (((_la-9)&-(0x1f+1)) == 0 && ((1<<uint((_la-9)))&((1<<(FqlParserOpenBracket-9))|(1<<(FqlParserOpenParen-9))|(1<<(FqlParserOpenBrace-9))|(1<<(FqlParserPlus-9))|(1<<(FqlParserMinus-9))|(1<<(FqlParserAnd-9))|(1<<(FqlParserOr-9))|(1<<(FqlParserFor-9))|(1<<(FqlParserReturn-9))|(1<<(FqlParserWaitfor-9))|(1<<(FqlParserDistinct-9))|(1<<(FqlParserFilter-9)))) != 0) || (((_la-41)&-(0x1f+1)) == 0 && ((1<<uint((_la-41)))&((1<<(FqlParserSort-41))|(1<<(FqlParserLimit-41))|(1<<(FqlParserLet-41))|(1<<(FqlParserCollect-41))|(1<<(FqlParserSortDirection-41))|(1<<(FqlParserNone-41))|(1<<(FqlParserNull-41))|(1<<(FqlParserBooleanLiteral-41))|(1<<(FqlParserUse-41))|(1<<(FqlParserInto-41))|(1<<(FqlParserKeep-41))|(1<<(FqlParserWith-41))|(1<<(FqlParserCount-41))|(1<<(FqlParserAll-41))|(1<<(FqlParserAny-41))|(1<<(FqlParserAggregate-41))|(1<<(FqlParserEvent-41))|(1<<(FqlParserLike-41))|(1<<(FqlParserNot-41))|(1<<(FqlParserIn-41))|(1<<(FqlParserParam-41))|(1<<(FqlParserIdentifier-41))|(1<<(FqlParserStringLiteral-41))|(1<<(FqlParserIntegerLiteral-41))|(1<<(FqlParserFloatLiteral-41))|(1<<(FqlParserNamespaceSegment-41)))) != 0) {
					{
						p.SetState(510)
						p.expression(0)
					}

				}
				{
					p.SetState(513)
					p.Match(FqlParserColon)
				}
				{
					p.SetState(514)
					p.expression(16)
				}

			}

		}
		p.SetState(519)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 46, p.GetParserRuleContext())
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

func (s *ExpressionGroupContext) ErrorOperator() IErrorOperatorContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IErrorOperatorContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IErrorOperatorContext)
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
	p.EnterRule(localctx, 94, FqlParserRULE_expressionGroup)

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
		p.SetState(520)
		p.Match(FqlParserOpenParen)
	}
	{
		p.SetState(521)
		p.expression(0)
	}
	{
		p.SetState(522)
		p.Match(FqlParserCloseParen)
	}
	p.SetState(524)
	p.GetErrorHandler().Sync(p)

	if p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 47, p.GetParserRuleContext()) == 1 {
		{
			p.SetState(523)
			p.ErrorOperator()
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
	localctx = NewMemberExpressionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 96, FqlParserRULE_memberExpression)

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
		p.SetState(526)
		p.MemberExpressionSource()
	}
	p.SetState(528)
	p.GetErrorHandler().Sync(p)
	_alt = 1
	for ok := true; ok; ok = _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		switch _alt {
		case 1:
			{
				p.SetState(527)
				p.MemberExpressionPath()
			}

		default:
			panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		}

		p.SetState(530)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 48, p.GetParserRuleContext())
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
	localctx = NewMemberExpressionSourceContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 98, FqlParserRULE_memberExpressionSource)

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

	p.SetState(537)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 49, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(532)
			p.Variable()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(533)
			p.Param()
		}

	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(534)
			p.ArrayLiteral()
		}

	case 4:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(535)
			p.ObjectLiteral()
		}

	case 5:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(536)
			p.FunctionCall()
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

func (s *FunctionCallContext) FunctionIdentifier() IFunctionIdentifierContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IFunctionIdentifierContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IFunctionIdentifierContext)
}

func (s *FunctionCallContext) Arguments() IArgumentsContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IArgumentsContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IArgumentsContext)
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
	localctx = NewFunctionCallContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 100, FqlParserRULE_functionCall)

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
		p.SetState(539)
		p.Namespace()
	}
	{
		p.SetState(540)
		p.FunctionIdentifier()
	}
	{
		p.SetState(541)
		p.Arguments()
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
	localctx = NewFunctionCallExpressionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 102, FqlParserRULE_functionCallExpression)

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
		p.SetState(543)
		p.FunctionCall()
	}
	p.SetState(545)
	p.GetErrorHandler().Sync(p)

	if p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 50, p.GetParserRuleContext()) == 1 {
		{
			p.SetState(544)
			p.ErrorOperator()
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

func (s *MemberExpressionPathContext) QuestionMark() antlr.TerminalNode {
	return s.GetToken(FqlParserQuestionMark, 0)
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
	localctx = NewMemberExpressionPathContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 104, FqlParserRULE_memberExpressionPath)
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

	p.SetState(557)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 53, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		p.SetState(548)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		if _la == FqlParserQuestionMark {
			{
				p.SetState(547)
				p.Match(FqlParserQuestionMark)
			}

		}
		{
			p.SetState(550)
			p.Match(FqlParserDot)
		}
		{
			p.SetState(551)
			p.PropertyName()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		p.SetState(554)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		if _la == FqlParserQuestionMark {
			{
				p.SetState(552)
				p.Match(FqlParserQuestionMark)
			}
			{
				p.SetState(553)
				p.Match(FqlParserDot)
			}

		}
		{
			p.SetState(556)
			p.ComputedPropertyName()
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
	localctx = NewErrorOperatorContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 106, FqlParserRULE_errorOperator)

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
		p.SetState(559)
		p.Match(FqlParserQuestionMark)
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
	p.EnterRule(localctx, 108, FqlParserRULE_functionIdentifier)
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
		p.SetState(561)
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
	p.EnterRule(localctx, 110, FqlParserRULE_namespace)
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
	p.SetState(566)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == FqlParserNamespaceSegment {
		{
			p.SetState(563)
			p.Match(FqlParserNamespaceSegment)
		}

		p.SetState(568)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
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
	p.EnterRule(localctx, 112, FqlParserRULE_arguments)
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
		p.SetState(569)
		p.Match(FqlParserOpenParen)
	}
	p.SetState(578)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if (((_la-9)&-(0x1f+1)) == 0 && ((1<<uint((_la-9)))&((1<<(FqlParserOpenBracket-9))|(1<<(FqlParserOpenParen-9))|(1<<(FqlParserOpenBrace-9))|(1<<(FqlParserPlus-9))|(1<<(FqlParserMinus-9))|(1<<(FqlParserAnd-9))|(1<<(FqlParserOr-9))|(1<<(FqlParserFor-9))|(1<<(FqlParserReturn-9))|(1<<(FqlParserWaitfor-9))|(1<<(FqlParserDistinct-9))|(1<<(FqlParserFilter-9)))) != 0) || (((_la-41)&-(0x1f+1)) == 0 && ((1<<uint((_la-41)))&((1<<(FqlParserSort-41))|(1<<(FqlParserLimit-41))|(1<<(FqlParserLet-41))|(1<<(FqlParserCollect-41))|(1<<(FqlParserSortDirection-41))|(1<<(FqlParserNone-41))|(1<<(FqlParserNull-41))|(1<<(FqlParserBooleanLiteral-41))|(1<<(FqlParserUse-41))|(1<<(FqlParserInto-41))|(1<<(FqlParserKeep-41))|(1<<(FqlParserWith-41))|(1<<(FqlParserCount-41))|(1<<(FqlParserAll-41))|(1<<(FqlParserAny-41))|(1<<(FqlParserAggregate-41))|(1<<(FqlParserEvent-41))|(1<<(FqlParserLike-41))|(1<<(FqlParserNot-41))|(1<<(FqlParserIn-41))|(1<<(FqlParserParam-41))|(1<<(FqlParserIdentifier-41))|(1<<(FqlParserStringLiteral-41))|(1<<(FqlParserIntegerLiteral-41))|(1<<(FqlParserFloatLiteral-41))|(1<<(FqlParserNamespaceSegment-41)))) != 0) {
		{
			p.SetState(570)
			p.expression(0)
		}
		p.SetState(575)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		for _la == FqlParserComma {
			{
				p.SetState(571)
				p.Match(FqlParserComma)
			}
			{
				p.SetState(572)
				p.expression(0)
			}

			p.SetState(577)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)
		}

	}
	{
		p.SetState(580)
		p.Match(FqlParserCloseParen)
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
	p.EnterRule(localctx, 114, FqlParserRULE_arrayOperator)
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
		p.SetState(582)
		_la = p.GetTokenStream().LA(1)

		if !(((_la-46)&-(0x1f+1)) == 0 && ((1<<uint((_la-46)))&((1<<(FqlParserNone-46))|(1<<(FqlParserAll-46))|(1<<(FqlParserAny-46)))) != 0) {
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
	p.EnterRule(localctx, 116, FqlParserRULE_inOperator)

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

	p.SetState(587)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case FqlParserIn:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(584)
			p.Match(FqlParserIn)
		}

	case FqlParserNot:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(585)
			p.Match(FqlParserNot)
		}
		{
			p.SetState(586)
			p.Match(FqlParserIn)
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
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
	localctx = NewLikeOperatorContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 118, FqlParserRULE_likeOperator)

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

	p.SetState(592)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case FqlParserLike:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(589)
			p.Match(FqlParserLike)
		}

	case FqlParserNot:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(590)
			p.Match(FqlParserNot)
		}
		{
			p.SetState(591)
			p.Match(FqlParserLike)
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
	p.EnterRule(localctx, 120, FqlParserRULE_equalityOperator)
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
		p.SetState(594)
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
	p.EnterRule(localctx, 122, FqlParserRULE_regexpOperator)
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
		p.SetState(596)
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
	p.EnterRule(localctx, 124, FqlParserRULE_logicalAndOperator)

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
		p.SetState(598)
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
	p.EnterRule(localctx, 126, FqlParserRULE_logicalOrOperator)

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
		p.SetState(600)
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
	p.EnterRule(localctx, 128, FqlParserRULE_multiplicativeOperator)
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
		p.SetState(602)
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
	p.EnterRule(localctx, 130, FqlParserRULE_additiveOperator)
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
		p.SetState(604)
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
	p.EnterRule(localctx, 132, FqlParserRULE_unaryOperator)
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
		p.SetState(606)
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
	p.EnterRule(localctx, 134, FqlParserRULE_param)

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
		p.SetState(608)
		p.Match(FqlParserParam)
	}
	{
		p.SetState(609)
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
	p.EnterRule(localctx, 136, FqlParserRULE_variable)

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
		p.SetState(611)
		p.Match(FqlParserIdentifier)
	}

	return localctx
}

func (p *FqlParser) Sempred(localctx antlr.RuleContext, ruleIndex, predIndex int) bool {
	switch ruleIndex {
	case 46:
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
		return p.Precpred(p.GetParserRuleContext(), 24)

	case 1:
		return p.Precpred(p.GetParserRuleContext(), 23)

	case 2:
		return p.Precpred(p.GetParserRuleContext(), 22)

	case 3:
		return p.Precpred(p.GetParserRuleContext(), 21)

	case 4:
		return p.Precpred(p.GetParserRuleContext(), 20)

	case 5:
		return p.Precpred(p.GetParserRuleContext(), 19)

	case 6:
		return p.Precpred(p.GetParserRuleContext(), 18)

	case 7:
		return p.Precpred(p.GetParserRuleContext(), 17)

	case 8:
		return p.Precpred(p.GetParserRuleContext(), 16)

	case 9:
		return p.Precpred(p.GetParserRuleContext(), 15)

	default:
		panic("No predicate with index: " + fmt.Sprint(predIndex))
	}
}
