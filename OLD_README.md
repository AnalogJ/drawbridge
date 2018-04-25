
> Sections below need to be re-written and added to README.md

----


# How does config file work:

- list of questions, each question has a key, that key is used to populate the template.
- each question has a description, used when asking the user for input
- each question has a type, a golang type, used when generating the struct, and for validaiton
- question can have validation, ensure that it's proper value
	- https://github.com/go-playground/validator
	- https://github.com/bluesuncorp/validator
	- https://github.com/xeipuuv/gojsonschema
	- https://github.com/thedevsaddam/govalidator
	- https://github.com/go-validator/validator
	- https://github.com/gima/govalid
	- https://github.com/lestrrat/go-jsref
	- https://medium.com/@lestrrat/json-schema-and-go-3c7439959077
	- https://github.com/lestrrat/go-jsschema

- question can have range of allowed values
- question can have an example string (not default), used for hinting to the user.
- question can have ui_group_by value, 1,2,3 used in ui for listing.
- question can have ui_hidden value, boolean, used in ui to hide during listing.

- questions will be used to create a dynamic Struct, with tags added dynamically: https://github.com/fatih/gomodifytags

- list of answers
- answers can reference an external file using `_file`, which will be loaded inplace.
- answers must provide atleast one of the questions. (empty objects will throw an error)
- answers will be validated against the questions. Any invalid answers removed? throw an error?

- template section
- custom/overridable templates supported:
	- config template
	- config filename template
	- ssh key filepath template
