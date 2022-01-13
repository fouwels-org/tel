// SPDX-FileCopyrightText: 2022 Kaelan Thijs Fouwels <kaelan.thijs@fouwels.com>
//
// SPDX-License-Identifier: MIT
package drivers

import "context"

type Driver interface {
	Run(ctx context.Context) error
}
