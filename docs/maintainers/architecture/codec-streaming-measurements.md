# Codec streaming measurements

Supporting data for the [migration report](codec-streaming-migration.md).
Every cell is the median of ten samples: **ns/op / B/op / allocs/op**.
Common encode/decode rows and initially noisy rows use separate serial repeats;
other rows use the complete 119-workload pass. See the report for repeat criteria.
All revisions use identical workload fixtures, Go 1.27.1, GOMAXPROCS=12,
`-benchmem -benchtime=200ms -count=10`, Darwin/arm64 on Apple M2 Max.
The final column is the uncommitted follow-up to d9b3ba82 described in the report.

| Workload | 35693eb3 | 87684ac7 | 852e08c7 | d9b3ba82 | final |
| --- | ---: | ---: | ---: | ---: | ---: |
| JSONCodecEncode/duration | 65.6 / 112 / 2 | 133.6 / 208 / 3 | 137.1 / 208 / 3 | 136.9 / 208 / 3 | 140.8 / 208 / 3 |
| JSONCodecEncode/regexp | 99.5 / 112 / 2 | 153 / 208 / 3 | 155.5 / 208 / 3 | 157.2 / 208 / 3 | 160.4 / 208 / 3 |
| JSONCodecEncode/range_1024 | 10,390 / 8,216 / 13 | 17,446 / 8,277 / 9 | 17,426 / 8,277 / 9 | 17,475 / 8,277 / 9 | 17,542 / 8,277 / 9 |
| JSONCodecEncode/flat_array_1024 | 26,086 / 11,152 / 932 | 24,176 / 8,277 / 9 | 24,150 / 8,277 / 9 | 24,204 / 8,277 / 9 | 26,307 / 8,277 / 9 |
| JSONCodecEncode/flat_object_256 | 16,676 / 12,867 / 422 | 17,399 / 12,457 / 267 | 17,487 / 12,457 / 267 | 17,526 / 12,457 / 267 | 18,436 / 12,457 / 267 |
| JSONCodecEncode/nested_array_10000 | 305,858 / 65,520 / 11 | 398,242 / 65,640 / 12 | 395,526 / 65,638 / 12 | 395,196 / 65,636 / 12 | 423,198 / 65,636 / 12 |
| JSONCodecEncode/nested_object_5000 | 758,504 / 630,523 / 15,011 | 767,850 / 550,825 / 15,012 | 762,622 / 550,826 / 15,012 | 762,082 / 550,827 / 15,012 | 826,390 / 550,822 / 15,012 |
| JSONCodecDecode/flat_array_1024 | 59,734 / 90,120 / 1,812 | 75,786 / 93,936 / 2,834 | 76,617 / 93,936 / 2,834 | 76,586 / 93,936 / 2,834 | 77,270 / 93,936 / 2,834 |
| JSONCodecDecode/flat_object_256 | 37,256 / 54,064 / 781 | 43,048 / 51,720 / 1,037 | 43,026 / 51,736 / 1,037 | 43,250 / 51,736 / 1,037 | 43,748 / 51,736 / 1,037 |
| JSONCodecDecode/nested_array_10000 | 775,828 / 2,307,975 / 20,029 | 986,555 / 2,609,337 / 20,053 | 955,352 / 2,609,339 / 20,053 | 949,782 / 2,609,337 / 20,053 | 980,445 / 2,609,338 / 20,053 |
| JSONCodecDecode/nested_object_5000 | 797,120 / 2,610,249 / 20,026 | 974,981 / 2,828,402 / 20,076 | 976,371 / 2,828,404 / 20,076 | 971,000 / 2,828,398 / 20,076 | 966,504 / 2,828,400 / 20,076 |
| MsgpackCodecEncode/regexp | 59.2 / 112 / 2 | 139 / 288 / 4 | 138.9 / 288 / 4 | 141.8 / 288 / 4 | 144.1 / 288 / 4 |
| MsgpackCodecEncode/range_1024 | 17,358 / 14,328 / 776 | 22,010 / 14,505 / 778 | 21,815 / 14,505 / 778 | 21,821 / 14,505 / 778 | 23,079 / 14,505 / 778 |
| MsgpackCodecEncode/flat_array_1024 | 12,830 / 8,181 / 8 | 17,466 / 8,357 / 10 | 16,872 / 8,357 / 10 | 16,820 / 8,357 / 10 | 18,422 / 8,357 / 10 |
| MsgpackCodecEncode/flat_object_256 | 11,054 / 8,261 / 264 | 13,378 / 8,437 / 266 | 13,316 / 8,437 / 266 | 13,177 / 8,437 / 266 | 13,796 / 8,437 / 266 |
| MsgpackCodecEncode/nested_array_10000 | 258,726 / 32,763 / 10 | 298,596 / 32,939 / 12 | 297,492 / 32,939 / 12 | 298,425 / 32,939 / 12 | 307,008 / 32,938 / 12 |
| MsgpackCodecEncode/nested_object_5000 | 644,854 / 512,958 / 10,010 | 694,684 / 513,134 / 10,012 | 688,734 / 513,136 / 10,012 | 693,212 / 513,134 / 10,012 | 700,904 / 513,130 / 10,012 |
| MsgpackCodecDecode/flat_array_1024 | 26,518 / 24,663 / 771 | 39,739 / 24,847 / 774 | 38,478 / 24,847 / 774 | 38,541 / 24,847 / 774 | 38,334 / 24,847 / 774 |
| MsgpackCodecDecode/flat_object_256 | 15,166 / 19,315 / 252 | 19,590 / 19,491 / 255 | 19,398 / 19,499 / 255 | 19,473 / 19,491 / 255 | 19,542 / 19,491 / 255 |
| MsgpackCodecDecode/nested_array_10000 | 506,624 / 400,212 / 20,001 | 559,296 / 400,389 / 20,004 | 555,589 / 400,394 / 20,004 | 551,896 / 400,390 / 20,004 | 552,884 / 400,392 / 20,004 |
| MsgpackCodecDecode/nested_object_5000 | 754,851 / 1,720,582 / 15,001 | 818,862 / 1,720,763 / 15,004 | 822,890 / 1,720,766 / 15,004 | 819,952 / 1,720,766 / 15,004 | 832,348 / 1,720,767 / 15,004 |
| CodecWriter/application/json/array/default/background/Buffer | 459,148 / 335,058 / 16,297 | 371,447 / 0 / 0 | 363,112 / 0 / 0 | 374,760 / 0 / 0 | 395,836 / 0 / 0 |
| CodecWriter/application/json/array/default/background/Discard | 457,660 / 335,059 / 16,297 | 355,063 / 0 / 0 | 357,645 / 0 / 0 | 358,224 / 0 / 0 | 375,374 / 0 / 0 |
| CodecWriter/application/json/array/default/cancellable/Buffer | 460,590 / 335,058 / 16,297 | 373,678 / 0 / 0 | 376,696 / 0 / 0 | 380,601 / 0 / 0 | 397,398 / 0 / 0 |
| CodecWriter/application/json/array/default/cancellable/Discard | 459,026 / 335,058 / 16,297 | 357,496 / 0 / 0 | 360,492 / 0 / 0 | 363,860 / 0 / 0 | 381,609 / 0 / 0 |
| CodecWriter/application/json/array/configured/background/Buffer | 470,348 / 335,058 / 16,297 | 379,756 / 0 / 0 | 386,418 / 0 / 0 | 390,043 / 0 / 0 | 410,634 / 0 / 0 |
| CodecWriter/application/json/array/configured/background/Discard | 469,893 / 335,058 / 16,297 | 362,838 / 0 / 0 | 368,788 / 0 / 0 | 372,765 / 0 / 0 | 394,938 / 0 / 0 |
| CodecWriter/application/json/array/configured/cancellable/Buffer | 470,816 / 335,058 / 16,297 | 385,065 / 0 / 0 | 394,166 / 0 / 0 | 391,702 / 0 / 0 | 411,362 / 0 / 0 |
| CodecWriter/application/json/array/configured/cancellable/Discard | 477,850 / 335,058 / 16,297 | 371,694 / 0 / 0 | 375,992 / 0 / 0 | 400,844 / 0 / 0 | 396,320 / 0 / 0 |
| CodecWriter/application/json/object/default/background/Buffer | 299,197 / 212,658 / 8,106 | 248,725 / 65,653 / 4,098 | 257,161 / 65,652 / 4,098 | 243,230 / 65,653 / 4,098 | 253,528 / 65,653 / 4,098 |
| CodecWriter/application/json/object/default/background/Discard | 293,656 / 212,658 / 8,106 | 241,820 / 65,653 / 4,098 | 250,380 / 65,653 / 4,098 | 235,144 / 65,654 / 4,098 | 244,600 / 65,653 / 4,098 |
| CodecWriter/application/json/object/default/cancellable/Buffer | 293,448 / 212,658 / 8,106 | 251,236 / 65,653 / 4,098 | 260,344 / 65,652 / 4,098 | 244,772 / 65,652 / 4,098 | 253,990 / 65,653 / 4,098 |
| CodecWriter/application/json/object/default/cancellable/Discard | 294,719 / 212,658 / 8,106 | 243,336 / 65,653 / 4,098 | 254,629 / 65,653 / 4,098 | 247,534 / 65,652 / 4,098 | 246,381 / 65,652 / 4,098 |
| CodecWriter/application/json/object/configured/background/Buffer | 282,532 / 212,654 / 8,106 | 253,046 / 65,653 / 4,098 | 263,446 / 65,652 / 4,098 | 256,612 / 65,652 / 4,098 | 257,779 / 65,652 / 4,098 |
| CodecWriter/application/json/object/configured/background/Discard | 302,728 / 212,658 / 8,106 | 244,944 / 65,653 / 4,098 | 257,436 / 65,653 / 4,098 | 247,776 / 65,652 / 4,098 | 250,844 / 65,653 / 4,098 |
| CodecWriter/application/json/object/configured/cancellable/Buffer | 297,428 / 212,658 / 8,106 | 253,782 / 65,652 / 4,098 | 248,612 / 65,652 / 4,098 | 258,034 / 65,652 / 4,098 | 258,638 / 65,652 / 4,098 |
| CodecWriter/application/json/object/configured/cancellable/Discard | 298,702 / 212,658 / 8,106 | 246,158 / 65,654 / 4,098 | 370,471 / 65,653 / 4,098 | 249,114 / 65,653 / 4,098 | 251,018 / 65,653 / 4,098 |
| CodecWriter/application/json/range/default/background/Buffer | 736,706 / 1,048,648 / 26 | 1,061,735 / 0 / 0 | 1,128,428 / 0 / 0 | 1,058,539 / 0 / 0 | 1,038,760 / 0 / 0 |
| CodecWriter/application/json/range/default/background/Discard | 725,696 / 1,048,648 / 26 | 981,334 / 0 / 0 | 958,596 / 0 / 0 | 976,618 / 0 / 0 | 958,611 / 0 / 0 |
| CodecWriter/application/json/range/default/cancellable/Buffer | 737,184 / 1,048,648 / 26 | 1,070,492 / 0 / 0 | 1,064,134 / 0 / 0 | 1,057,746 / 0 / 0 | 1,043,056 / 0 / 0 |
| CodecWriter/application/json/range/default/cancellable/Discard | 725,058 / 1,048,648 / 26 | 1,005,832 / 0 / 0 | 1,000,105 / 0 / 0 | 996,750 / 0 / 0 | 981,682 / 0 / 0 |
| CodecWriter/application/json/range/configured/background/Buffer | 2,114,232 / 2,093,544 / 130,731 | 1,939,341 / 522,490 / 65,280 | 1,944,804 / 522,492 / 65,280 | 1,899,358 / 522,484 / 65,280 | 2,001,440 / 522,496 / 65,280 |
| CodecWriter/application/json/range/configured/background/Discard | 2,117,910 / 2,093,543 / 130,731 | 1,881,316 / 522,525 / 65,280 | 1,871,692 / 522,524 / 65,280 | 1,836,759 / 522,520 / 65,280 | 1,950,448 / 522,522 / 65,280 |
| CodecWriter/application/json/range/configured/cancellable/Buffer | 2,122,440 / 2,093,544 / 130,731 | 1,944,783 / 522,492 / 65,280 | 1,969,868 / 522,493 / 65,280 | 1,902,812 / 522,484 / 65,280 | 2,115,668 / 522,495 / 65,280 |
| CodecWriter/application/json/range/configured/cancellable/Discard | 2,117,847 / 2,093,543 / 130,731 | 1,884,526 / 522,526 / 65,280 | 1,877,286 / 522,525 / 65,280 | 1,845,181 / 522,521 / 65,280 | 2,050,482 / 522,530 / 65,280 |
| CodecWriter/application/json/string_1MiB/default/background/Buffer | 4,302,054 / 4,194,288 / 17 | 2,777,905 / 0 / 0 | 2,844,340 / 0 / 0 | 2,800,919 / 0 / 0 | 2,873,030 / 0 / 0 |
| CodecWriter/application/json/string_1MiB/default/background/Discard | 4,296,815 / 4,194,288 / 17 | 2,740,666 / 0 / 0 | 2,576,968 / 0 / 0 | 2,781,836 / 0 / 0 | 2,853,845 / 0 / 0 |
| CodecWriter/application/json/string_1MiB/default/cancellable/Buffer | 4,292,537 / 4,194,288 / 17 | 2,772,004 / 0 / 0 | 2,597,348 / 0 / 0 | 2,798,344 / 0 / 0 | 2,884,603 / 0 / 0 |
| CodecWriter/application/json/string_1MiB/default/cancellable/Discard | 4,295,156 / 4,194,288 / 17 | 2,738,795 / 0 / 0 | 2,575,946 / 0 / 0 | 2,769,504 / 0 / 0 | 2,862,384 / 0 / 0 |
| CodecWriter/application/json/string_1MiB/configured/background/Buffer | 4,314,728 / 4,194,288 / 17 | 2,773,388 / 0 / 0 | 2,873,830 / 0 / 0 | 2,778,608 / 0 / 0 | 2,888,461 / 0 / 0 |
| CodecWriter/application/json/string_1MiB/configured/background/Discard | 4,303,088 / 4,194,288 / 17 | 2,740,276 / 0 / 0 | 2,892,422 / 0 / 0 | 2,757,698 / 0 / 0 | 2,864,416 / 0 / 0 |
| CodecWriter/application/json/string_1MiB/configured/cancellable/Buffer | 4,288,807 / 4,194,288 / 17 | 2,748,548 / 0 / 0 | 2,595,833 / 0 / 0 | 2,779,090 / 0 / 0 | 2,803,140 / 0 / 0 |
| CodecWriter/application/json/string_1MiB/configured/cancellable/Discard | 4,294,206 / 4,194,288 / 17 | 2,725,152 / 0 / 0 | 2,576,932 / 0 / 0 | 2,756,242 / 0 / 0 | 2,683,194 / 0 / 0 |
| CodecWriter/application/json/binary_1MiB/default/background/Buffer | 6,261,676 / 6,995,953 / 19 | 467,990 / 0 / 0 | 461,172 / 0 / 0 | 462,588 / 0 / 0 | 461,673 / 0 / 0 |
| CodecWriter/application/json/binary_1MiB/default/background/Discard | 6,175,754 / 6,995,952 / 19 | 440,814 / 0 / 0 | 434,791 / 0 / 0 | 435,463 / 0 / 0 | 434,550 / 0 / 0 |
| CodecWriter/application/json/binary_1MiB/default/cancellable/Buffer | 6,220,368 / 6,995,953 / 19 | 470,117 / 0 / 0 | 461,912 / 0 / 0 | 463,787 / 0 / 0 | 461,118 / 0 / 0 |
| CodecWriter/application/json/binary_1MiB/default/cancellable/Discard | 6,122,070 / 6,995,952 / 19 | 449,324 / 0 / 0 | 442,652 / 0 / 0 | 443,172 / 0 / 0 | 442,999 / 0 / 0 |
| CodecWriter/application/json/binary_1MiB/configured/background/Buffer | 6,233,238 / 6,995,953 / 19 | 468,907 / 0 / 0 | 462,304 / 0 / 0 | 462,620 / 0 / 0 | 459,886 / 0 / 0 |
| CodecWriter/application/json/binary_1MiB/configured/background/Discard | 6,145,556 / 6,995,952 / 19 | 440,894 / 0 / 0 | 439,612 / 0 / 0 | 436,216 / 0 / 0 | 433,680 / 0 / 0 |
| CodecWriter/application/json/binary_1MiB/configured/cancellable/Buffer | 6,235,044 / 6,995,952 / 19 | 471,216 / 0 / 0 | 462,979 / 0 / 0 | 463,039 / 0 / 0 | 461,132 / 0 / 0 |
| CodecWriter/application/json/binary_1MiB/configured/cancellable/Discard | 6,184,482 / 6,995,952 / 19 | 450,529 / 0 / 0 | 442,880 / 0 / 0 | 443,012 / 0 / 0 | 442,058 / 0 / 0 |
| CodecWriter/application/vnd.msgpack/array/default/background/Buffer | 208,044 / 131,121 / 12 | 256,788 / 80 / 1 | 251,402 / 80 / 1 | 251,006 / 80 / 1 | 281,383 / 80 / 1 |
| CodecWriter/application/vnd.msgpack/array/default/background/Discard | 207,647 / 131,122 / 12 | 218,194 / 80 / 1 | 214,080 / 80 / 1 | 217,757 / 80 / 1 | 241,810 / 80 / 1 |
| CodecWriter/application/vnd.msgpack/array/default/cancellable/Buffer | 208,228 / 131,121 / 12 | 265,978 / 80 / 1 | 251,868 / 80 / 1 | 256,706 / 80 / 1 | 276,246 / 80 / 1 |
| CodecWriter/application/vnd.msgpack/array/default/cancellable/Discard | 207,826 / 131,122 / 12 | 232,082 / 80 / 1 | 218,714 / 80 / 1 | 223,380 / 80 / 1 | 242,632 / 80 / 1 |
| CodecWriter/application/vnd.msgpack/array/configured/background/Buffer | 233,236 / 131,122 / 12 | 279,738 / 80 / 1 | 275,352 / 80 / 1 | 280,072 / 80 / 1 | 304,293 / 80 / 1 |
| CodecWriter/application/vnd.msgpack/array/configured/background/Discard | 231,502 / 131,122 / 12 | 247,192 / 80 / 1 | 242,160 / 80 / 1 | 246,748 / 80 / 1 | 271,091 / 80 / 1 |
| CodecWriter/application/vnd.msgpack/array/configured/cancellable/Buffer | 231,578 / 131,122 / 12 | 289,996 / 80 / 1 | 280,962 / 80 / 1 | 281,274 / 80 / 1 | 305,200 / 80 / 1 |
| CodecWriter/application/vnd.msgpack/array/configured/cancellable/Discard | 231,834 / 131,122 / 12 | 250,890 / 80 / 1 | 247,420 / 80 / 1 | 247,747 / 80 / 1 | 271,940 / 80 / 1 |
| CodecWriter/application/vnd.msgpack/object/default/background/Buffer | 168,697 / 131,205 / 4,108 | 186,340 / 65,732 / 4,098 | 186,791 / 65,731 / 4,098 | 186,688 / 65,731 / 4,098 | 209,510 / 65,732 / 4,098 |
| CodecWriter/application/vnd.msgpack/object/default/background/Discard | 168,326 / 131,205 / 4,108 | 164,921 / 65,732 / 4,098 | 166,027 / 65,732 / 4,098 | 166,654 / 65,732 / 4,098 | 179,261 / 65,732 / 4,098 |
| CodecWriter/application/vnd.msgpack/object/default/cancellable/Buffer | 168,342 / 131,204 / 4,108 | 187,681 / 65,732 / 4,098 | 188,222 / 65,731 / 4,098 | 187,978 / 65,731 / 4,098 | 212,130 / 65,731 / 4,098 |
| CodecWriter/application/vnd.msgpack/object/default/cancellable/Discard | 168,564 / 131,205 / 4,108 | 171,994 / 65,732 / 4,098 | 169,412 / 65,731 / 4,098 | 169,204 / 65,731 / 4,098 | 186,268 / 65,732 / 4,098 |
| CodecWriter/application/vnd.msgpack/object/configured/background/Buffer | 175,330 / 131,204 / 4,108 | 197,553 / 65,731 / 4,098 | 193,102 / 65,731 / 4,098 | 193,661 / 65,731 / 4,098 | 216,873 / 65,731 / 4,098 |
| CodecWriter/application/vnd.msgpack/object/configured/background/Discard | 174,614 / 131,206 / 4,108 | 176,091 / 65,731 / 4,098 | 173,116 / 65,731 / 4,098 | 173,457 / 65,731 / 4,098 | 188,248 / 65,732 / 4,098 |
| CodecWriter/application/vnd.msgpack/object/configured/cancellable/Buffer | 174,780 / 131,205 / 4,108 | 199,794 / 65,732 / 4,098 | 195,945 / 65,731 / 4,098 | 196,998 / 65,731 / 4,098 | 204,132 / 65,732 / 4,098 |
| CodecWriter/application/vnd.msgpack/object/configured/cancellable/Discard | 174,579 / 131,205 / 4,108 | 180,772 / 65,732 / 4,098 | 177,716 / 65,732 / 4,098 | 177,612 / 65,732 / 4,098 | 201,331 / 65,732 / 4,098 |
| CodecWriter/application/vnd.msgpack/range/default/background/Buffer | 1,175,248 / 1,047,015 / 65,294 | 1,379,476 / 522,580 / 65,281 | 1,349,659 / 522,587 / 65,281 | 1,369,861 / 522,578 / 65,281 | 1,444,449 / 522,585 / 65,281 |
| CodecWriter/application/vnd.msgpack/range/default/background/Discard | 1,172,600 / 1,047,030 / 65,294 | 1,214,896 / 522,607 / 65,281 | 1,219,538 / 522,598 / 65,281 | 1,235,413 / 522,600 / 65,281 | 1,373,283 / 522,600 / 65,281 |
| CodecWriter/application/vnd.msgpack/range/default/cancellable/Buffer | 1,181,277 / 1,047,014 / 65,294 | 1,453,930 / 522,584 / 65,281 | 1,374,716 / 522,580 / 65,281 | 1,373,896 / 522,580 / 65,281 | 1,511,760 / 522,584 / 65,281 |
| CodecWriter/application/vnd.msgpack/range/default/cancellable/Discard | 1,169,792 / 1,047,027 / 65,294 | 1,318,508 / 522,600 / 65,281 | 1,240,220 / 522,602 / 65,281 | 1,240,237 / 522,602 / 65,281 | 1,374,368 / 522,600 / 65,281 |
| CodecWriter/application/vnd.msgpack/range/configured/background/Buffer | 1,279,772 / 1,047,015 / 65,294 | 1,559,500 / 522,580 / 65,281 | 1,469,148 / 522,588 / 65,281 | 1,470,309 / 522,586 / 65,281 | 1,629,371 / 522,580 / 65,281 |
| CodecWriter/application/vnd.msgpack/range/configured/background/Discard | 1,271,668 / 1,047,030 / 65,294 | 1,399,674 / 522,596 / 65,281 | 1,316,494 / 522,600 / 65,281 | 1,335,171 / 522,595 / 65,281 | 1,500,843 / 522,602 / 65,281 |
| CodecWriter/application/vnd.msgpack/range/configured/cancellable/Buffer | 1,276,641 / 1,047,016 / 65,294 | 1,544,272 / 522,580 / 65,281 | 1,454,676 / 522,585 / 65,281 | 1,470,946 / 522,582 / 65,281 | 1,642,874 / 522,580 / 65,281 |
| CodecWriter/application/vnd.msgpack/range/configured/cancellable/Discard | 1,219,394 / 1,047,032 / 65,294 | 1,390,800 / 522,598 / 65,281 | 1,325,414 / 522,601 / 65,281 | 1,339,692 / 522,594 / 65,281 | 1,436,055 / 522,605 / 65,281 |
| CodecWriter/application/vnd.msgpack/string_1MiB/default/background/Buffer | 89,797 / 1,057,183 / 3 | 24,982 / 80 / 1 | 25,480 / 80 / 1 | 23,890 / 80 / 1 | 23,970 / 80 / 1 |
| CodecWriter/application/vnd.msgpack/string_1MiB/default/background/Discard | 86,156 / 1,057,240 / 3 | 1,014 / 80 / 1 | 996.3 / 80 / 1 | 996.2 / 80 / 1 | 1,008 / 80 / 1 |
| CodecWriter/application/vnd.msgpack/string_1MiB/default/cancellable/Buffer | 89,269 / 1,057,183 / 3 | 25,460 / 80 / 1 | 26,166 / 80 / 1 | 25,358 / 80 / 1 | 26,591 / 80 / 1 |
| CodecWriter/application/vnd.msgpack/string_1MiB/default/cancellable/Discard | 82,900 / 1,057,206 / 3 | 2,012 / 80 / 1 | 1,973 / 80 / 1 | 1,974 / 80 / 1 | 2,019 / 80 / 1 |
| CodecWriter/application/vnd.msgpack/string_1MiB/configured/background/Buffer | 95,956 / 1,057,182 / 3 | 25,171 / 80 / 1 | 25,872 / 80 / 1 | 24,578 / 80 / 1 | 25,112 / 80 / 1 |
| CodecWriter/application/vnd.msgpack/string_1MiB/configured/background/Discard | 69,862 / 1,057,094 / 3 | 1,030 / 80 / 1 | 1,010 / 80 / 1 | 998.7 / 80 / 1 | 1,044 / 80 / 1 |
| CodecWriter/application/vnd.msgpack/string_1MiB/configured/cancellable/Buffer | 92,567 / 1,057,240 / 3 | 25,692 / 80 / 1 | 26,106 / 80 / 1 | 25,224 / 80 / 1 | 25,836 / 80 / 1 |
| CodecWriter/application/vnd.msgpack/string_1MiB/configured/cancellable/Discard | 65,178 / 1,057,056 / 3 | 2,009 / 80 / 1 | 1,972 / 80 / 1 | 1,975 / 80 / 1 | 2,012 / 80 / 1 |
| CodecWriter/application/vnd.msgpack/binary_1MiB/default/background/Buffer | 90,669 / 1,057,240 / 3 | 24,769 / 80 / 1 | 24,532 / 80 / 1 | 23,926 / 80 / 1 | 23,978 / 80 / 1 |
| CodecWriter/application/vnd.msgpack/binary_1MiB/default/background/Discard | 65,228 / 1,057,054 / 3 | 1,032 / 80 / 1 | 992.5 / 80 / 1 | 996 / 80 / 1 | 1,022 / 80 / 1 |
| CodecWriter/application/vnd.msgpack/binary_1MiB/default/cancellable/Buffer | 91,999 / 1,057,240 / 3 | 25,155 / 80 / 1 | 26,914 / 80 / 1 | 25,165 / 80 / 1 | 25,700 / 80 / 1 |
| CodecWriter/application/vnd.msgpack/binary_1MiB/default/cancellable/Discard | 64,966 / 1,057,061 / 3 | 2,106 / 80 / 1 | 1,969 / 80 / 1 | 2,016 / 80 / 1 | 2,004 / 80 / 1 |
| CodecWriter/application/vnd.msgpack/binary_1MiB/configured/background/Buffer | 90,805 / 1,057,239 / 3 | 26,066 / 80 / 1 | 24,545 / 80 / 1 | 24,616 / 80 / 1 | 24,906 / 80 / 1 |
| CodecWriter/application/vnd.msgpack/binary_1MiB/configured/background/Discard | 65,002 / 1,057,058 / 3 | 1,046 / 80 / 1 | 1,008 / 80 / 1 | 999.6 / 80 / 1 | 1,020 / 80 / 1 |
| CodecWriter/application/vnd.msgpack/binary_1MiB/configured/cancellable/Buffer | 92,760 / 1,057,240 / 3 | 25,678 / 80 / 1 | 26,912 / 80 / 1 | 25,357 / 80 / 1 | 25,305 / 80 / 1 |
| CodecWriter/application/vnd.msgpack/binary_1MiB/configured/cancellable/Discard | 64,248 / 1,057,042 / 3 | 1,972 / 80 / 1 | 1,975 / 80 / 1 | 1,978 / 80 / 1 | 1,980 / 80 / 1 |
| CodecIterable/application/json/owned/Buffer | 43,706 / 58,441 / 1,709 | 31,960 / 6,148 / 768 | 31,932 / 6,148 / 768 | 31,834 / 6,148 / 768 | 32,116 / 6,148 / 768 |
| CodecIterable/application/json/owned/Discard | 40,728 / 58,440 / 1,709 | 30,616 / 6,147 / 768 | 30,532 / 6,147 / 768 | 30,476 / 6,147 / 768 | 30,808 / 6,147 / 768 |
| CodecIterable/application/json/borrowed/Buffer | 40,738 / 58,441 / 1,709 | 32,020 / 6,148 / 768 | 31,884 / 6,148 / 768 | 31,846 / 6,148 / 768 | 32,129 / 6,148 / 768 |
| CodecIterable/application/json/borrowed/Discard | 42,206 / 58,440 / 1,709 | 30,640 / 6,147 / 768 | 30,473 / 6,147 / 768 | 30,418 / 6,147 / 768 | 30,815 / 6,147 / 768 |
| CodecIterable/application/vnd.msgpack/owned/Buffer | 24,560 / 55,493 / 785 | 31,386 / 66,044 / 782 | 31,490 / 66,044 / 782 | 31,464 / 66,044 / 782 | 32,312 / 66,044 / 782 |
| CodecIterable/application/vnd.msgpack/owned/Discard | 27,478 / 55,494 / 785 | 28,885 / 66,042 / 782 | 29,504 / 66,043 / 782 | 29,628 / 66,043 / 782 | 30,577 / 66,043 / 782 |
| CodecIterable/application/vnd.msgpack/borrowed/Buffer | 27,106 / 55,494 / 785 | 30,714 / 66,043 / 782 | 31,448 / 66,044 / 782 | 31,438 / 66,044 / 782 | 32,228 / 66,044 / 782 |
| CodecIterable/application/vnd.msgpack/borrowed/Discard | 26,707 / 55,494 / 785 | 28,793 / 66,042 / 782 | 29,628 / 66,043 / 782 | 29,715 / 66,043 / 782 | 30,593 / 66,043 / 782 |
| SessionCodecIterable/application/json/scalars | 52,159 / 59,609 / 1,729 | 39,516 / 15,854 / 802 | 39,510 / 15,854 / 802 | 39,656 / 15,982 / 805 | 51,329 / 15,718 / 801 |
| SessionCodecIterable/application/json/resources | 300,486 / 342,651 / 3,913 | 288,788 / 302,065 / 3,911 | 290,618 / 302,071 / 3,911 | 290,762 / 302,197 / 3,914 | 323,175 / 301,929 / 3,910 |
| SessionCodecIterable/application/json/containers/depth_1 | 487,262 / 498,426 / 6,987 | 505,952 / 441,516 / 6,984 | 471,144 / 441,486 / 6,984 | 475,634 / 441,617 / 6,987 | 652,570 / 441,618 / 6,986 |
| SessionCodecIterable/application/json/containers/depth_8 | 1,857,514 / 1,498,282 / 28,494 | 1,827,798 / 1,326,777 / 28,491 | 1,773,731 / 1,326,760 / 28,491 | 1,771,252 / 1,326,896 / 28,494 | 2,654,482 / 1,328,284 / 28,498 |
| SessionCodecIterable/application/json/containers/depth_32 | 6,975,486 / 5,062,154 / 102,225 | 6,930,072 / 4,497,544 / 102,222 | 6,786,154 / 4,497,494 / 102,222 | 7,056,552 / 4,497,724 / 102,226 | 10,193,488 / 4,504,510 / 102,234 |
| SessionCodecIterable/application/vnd.msgpack/scalars | 38,202 / 56,680 / 805 | 49,204 / 75,769 / 816 | 48,144 / 75,766 / 816 | 52,030 / 75,920 / 820 | 56,092 / 75,633 / 815 |
| SessionCodecIterable/application/vnd.msgpack/resources | 258,184 / 416,432 / 3,913 | 280,976 / 435,514 / 3,924 | 266,097 / 435,504 / 3,924 | 295,120 / 435,634 / 3,928 | 295,136 / 435,344 / 3,923 |
| SessionCodecIterable/application/vnd.msgpack/containers/depth_1 | 407,296 / 531,185 / 5,962 | 433,348 / 550,264 / 5,973 | 424,326 / 550,258 / 5,973 | 551,682 / 550,624 / 5,980 | 583,118 / 550,333 / 5,975 |
| SessionCodecIterable/application/vnd.msgpack/containers/depth_8 | 1,493,001 / 1,448,960 / 20,303 | 1,586,432 / 1,467,995 / 20,314 | 1,572,497 / 1,468,004 / 20,314 | 2,304,325 / 1,469,762 / 20,325 | 2,395,440 / 1,469,469 / 20,320 |
| SessionCodecIterable/application/vnd.msgpack/containers/depth_32 | 5,598,496 / 4,594,822 / 69,458 | 5,927,542 / 4,613,848 / 69,468 | 5,841,788 / 4,613,898 / 69,469 | 8,867,428 / 4,621,040 / 69,484 | 9,106,469 / 4,620,781 / 69,480 |

## Adoption controls

| Workload | d9b3ba82 | final |
| --- | ---: | ---: |
| CodecAdoption/application/json/resources/absent | 25,349 / 8,292 / 1,026 | 25,408 / 8,292 / 1,026 |
| CodecAdoption/application/json/resources/noop | 25,339 / 8,292 / 1,026 | 27,598 / 8,316 / 1,027 |
| CodecAdoption/application/json/resources/result | 25,390 / 8,420 / 1,029 | 154,391 / 228,636 / 1,064 |
| CodecAdoption/application/json/containers_depth_8/absent | 199,421 / 8,294 / 1,026 | 214,338 / 8,292 / 1,026 |
| CodecAdoption/application/json/containers_depth_8/noop | 199,038 / 8,294 / 1,026 | 223,640 / 8,317 / 1,027 |
| CodecAdoption/application/json/containers_depth_8/result | 199,454 / 8,422 / 1,029 | 649,144 / 228,928 / 1,067 |
| CodecAdoption/application/vnd.msgpack/resources/absent | 44,538 / 68,188 / 1,040 | 49,358 / 68,188 / 1,040 |
| CodecAdoption/application/vnd.msgpack/resources/noop | 45,073 / 68,212 / 1,041 | 51,012 / 68,212 / 1,041 |
| CodecAdoption/application/vnd.msgpack/resources/result | 142,208 / 288,552 / 1,078 | 175,578 / 288,534 / 1,078 |
| CodecAdoption/application/vnd.msgpack/containers_depth_8/absent | 197,928 / 68,192 / 1,040 | 214,640 / 68,192 / 1,040 |
| CodecAdoption/application/vnd.msgpack/containers_depth_8/noop | 198,207 / 68,216 / 1,041 | 221,030 / 68,216 / 1,041 |
| CodecAdoption/application/vnd.msgpack/containers_depth_8/result | 510,257 / 288,837 / 1,081 | 579,456 / 288,822 / 1,081 |

## Callback transport control

The temporary field variant contains the final ownership and validation fixes
but replaces the retained private context cache with a callback field in Operation.
The variants change only callback transport; production retains the wrapper under
the approved performance-gate exception.

| Workload | retained wrapper | callback field |
| --- | ---: | ---: |
| JSONCodecEncode/duration | 140.8 / 208 / 3 | 147.8 / 208 / 3 |
| JSONCodecEncode/regexp | 160.4 / 208 / 3 | 167.1 / 208 / 3 |
| JSONCodecEncode/range_1024 | 17,542 / 8,277 / 9 | 17,481 / 8,277 / 9 |
| JSONCodecEncode/flat_array_1024 | 26,307 / 8,277 / 9 | 26,054 / 8,277 / 9 |
| JSONCodecEncode/flat_object_256 | 18,436 / 12,457 / 267 | 18,142 / 12,457 / 267 |
| JSONCodecEncode/nested_array_10000 | 423,198 / 65,636 / 12 | 428,569 / 65,638 / 12 |
| JSONCodecEncode/nested_object_5000 | 826,390 / 550,822 / 15,012 | 822,972 / 550,815 / 15,012 |
| JSONCodecDecode/flat_array_1024 | 77,270 / 93,936 / 2,834 | 76,238 / 93,952 / 2,834 |
| JSONCodecDecode/flat_object_256 | 43,748 / 51,736 / 1,037 | 43,486 / 51,752 / 1,037 |
| JSONCodecDecode/nested_array_10000 | 980,445 / 2,609,338 / 20,053 | 966,460 / 2,609,354 / 20,053 |
| JSONCodecDecode/nested_object_5000 | 966,504 / 2,828,400 / 20,076 | 974,096 / 2,828,415 / 20,076 |
| MsgpackCodecEncode/regexp | 144.1 / 288 / 4 | 159.1 / 304 / 4 |
| MsgpackCodecEncode/range_1024 | 23,079 / 14,505 / 778 | 23,371 / 14,521 / 778 |
| MsgpackCodecEncode/flat_array_1024 | 18,422 / 8,357 / 10 | 18,670 / 8,373 / 10 |
| MsgpackCodecEncode/flat_object_256 | 13,796 / 8,437 / 266 | 13,945 / 8,453 / 266 |
| MsgpackCodecEncode/nested_array_10000 | 307,008 / 32,938 / 12 | 308,844 / 32,953 / 12 |
| MsgpackCodecEncode/nested_object_5000 | 700,904 / 513,130 / 10,012 | 706,634 / 513,150 / 10,012 |
| MsgpackCodecDecode/flat_array_1024 | 38,334 / 24,847 / 774 | 38,500 / 24,863 / 774 |
| MsgpackCodecDecode/flat_object_256 | 19,542 / 19,491 / 255 | 19,518 / 19,507 / 255 |
| MsgpackCodecDecode/nested_array_10000 | 552,884 / 400,392 / 20,004 | 553,021 / 400,405 / 20,004 |
| MsgpackCodecDecode/nested_object_5000 | 832,348 / 1,720,767 / 15,004 | 821,080 / 1,720,784 / 15,004 |
