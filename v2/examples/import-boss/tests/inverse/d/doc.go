// c imports non-prod code. It shouldn't.
package d

import (
	_ "k8s.io/gengo/v2/examples/import-boss/tests/inverse/lib/nonprod"
)
