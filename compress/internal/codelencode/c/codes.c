/*
 * Copyright (c) 2023, Intel Corporation.
 * SPDX-License-Identifier: BSD-3-Clause
 */

#include <inttypes.h>


struct Histogram
{
	int32_t LiteralCodes[286]; // Literal Codes
	uint32_t re[2];
	int32_t DistanceCodes[30]; // Distance Codes
	uint32_t re2[2];
};


void prepareForCodeLenCode(struct Histogram *hist, uint8_t *source,uint16_t * lf,uint16_t * df)
{
	uint16_t lit_num = 0;
	uint16_t dis_num = 0;

	for (int32_t i = 285; i >= 0; i--)
	{
		if (hist->LiteralCodes[i] != 0)
		{
			lit_num = i + 1;
			break;
		}
	}
	for (int32_t i = 29; i >= 0; i--)
	{
		if (hist->DistanceCodes[i] != 0)
		{
			dis_num = i + 1;
			break;
		}
	}
	uint64_t insertOneDistance = 0;

	if (dis_num == 0)
	{
		dis_num = 1;
		insertOneDistance = 1;
	}

	// prepare to combine repeat numbers
	for (uint16_t i = (0); i < lit_num; i++)
	{
		source[i] = (uint8_t)(hist->LiteralCodes[i] >> 15);
	}

	for (uint16_t i = (0); i < dis_num; i++)
	{
		source[lit_num + i] = (uint8_t)(hist->DistanceCodes[i] >> 15);
	}
	if (insertOneDistance)
	{
		source[lit_num] = 1;
	}

	*lf = lit_num;
	*df = dis_num;
}
