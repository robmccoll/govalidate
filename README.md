govalidate
==========

Validate struct fields (currently strings and []strings) in Go using tags.

Define your struct with the validation built-in:

    type Valid struct {
        Required  string   `require:"+"`
        Number    string   `goodChars:"0123456789"`
        NotNumber string   `badChars:"0123456789"`
        Match     string   `match:"[0123456789]+"`
        Thing     []string `allowVals:"thing,thing2"`
    }

This keeps your validation and structure from getting out of sync.  You can apply
one of the tags up to all of the tags to each field.

Validate an instance of the struct by calling Validate:
    
    var inst Valid
    err := Validate(inst)

    if err != nil {
      fmt.Println(err.Error())
    }

Validate will return an error string describing which field value failed and why. Optionally,
you can call ```ValidateUseName(inst, "json")``` with ```"json"``` or some other tag name to use the values
in that tag in error messages in place of field names.

Validate tags:
- allowVals (ignored if empty string) is a comma-separated list of strings that are the only
  valid values for the givne string.
- goodChars (ignored if empty string) is a string of characters that are the only valid
  characters that may be in the string.
- badChars (ignored if empty string) is a string of characters that are explicitly disallowed
  from being in the string.
- match (ignored if empty string) is a regular expression that the given string must match
- require (ignored if empty string) if this contains any string at all, the length of the field
  must be greater than zero
