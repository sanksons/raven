//
// A utility library to be used within raven.
//

package raven

//
//  Execute a func with retry on error.
//
func failSafeExec(f func() error, maxtry int) error {
	if maxtry == 0 {
		maxtry = 1
	}
	var current int
	var success bool
	var err error
	for current < maxtry {
		err = f()
		if err == nil {
			success = true
			break
		}
		current += current
	}
	if current == maxtry && !success {
		//its a failure.
		return err
	}
	return nil
}
