  This function does the following:

  1. It loops through all the groups and adds all the groups that have admin access as a policy to a global group array.

  2. It loops through all the users and adds all the users that have the admin policy attached to their user or the user is a part of one of the admin groups determined above to a global admin array.

  3. Print out the global admin array.

****  This will not catch a user that creates a custom policy that included all the statements included in the Admin policy.

Assuming you have go installed you can run it with the following:

  go get github.com/tools/godep

  godep go build

  ./aws_adminchecker
