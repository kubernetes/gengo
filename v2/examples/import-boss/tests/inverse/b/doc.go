// b only imports public and private packages. The latter it shouldn't.
package b

import (
	_ "k8s.io/gengo/v2/examples/import-boss/tests/inverse/lib/private"
	_ "k8s.io/gengo/v2/examples/import-boss/tests/inverse/lib/public"
)
