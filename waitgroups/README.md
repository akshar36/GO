wg.Add(1) - means we are adding a goRoutine to the waitGroup.

wg.Done() - indicates that the goRoutine which is added to the waitgroup is completed (Basically removing it from the waitGroup).

defer wg.Done() - executes the wg.Done() finallly (After all statements in the function).

wg.Wait() - waits for all the goRoutines to be completed(waitgroup has no goRoutines in it).

waitGroups should be passed by reference as Go by default is a pass by value language you don't want to end up creating new copies of waitGroups and up manipulating it.

