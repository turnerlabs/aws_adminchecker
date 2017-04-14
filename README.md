  This function does the following:
  1. Get all the groups that have admin access added as a policy and add those groups to a global array
  2. Loop through all the users in the default aws account(via saml tool hopefully)

  a. If the user has the admin policy attached to their user, add them to the admin array.

  b. If the user is a part of one of the admin groups determined above then add them to the admin array.
  3. Print out the global admin user array.

****  This will not catch a user that creates a custom policy that included all the statements included in the Admin policy.

Assuming you have go installed you can run it with the following:

  go get github.com/tools/godep

  godep go build

  ./aws_adminchecker
