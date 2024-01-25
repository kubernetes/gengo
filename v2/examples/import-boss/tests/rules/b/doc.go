// b only public and private packages. The latter it shouldn't.
package b

import (
	_ "k8s.io/gengo/v2/examples/import-boss/tests/rules/c"
)
