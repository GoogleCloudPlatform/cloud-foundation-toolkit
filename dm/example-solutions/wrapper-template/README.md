# Wrapper Templates

Using wrapper templates is a clean way to extend or restrict already existing templates without modifying them.
This is a common practice when using templates from external sources for example from the Cloud Foundation Toolkit.

## Flexible solutions

Python provides easy access to the properties which will be passed forward to the GCP APIs. A wrapper template
is a good place to manipulate properties, for example to enforce naming conventions.

## Naming convention - Folders wrapper

In the *folders-wrapper.py* at line 11-12 the template is modifying the *Display Name* of the folders by adding a prefix.
This simple example can be easily extended, the prefix can be loaded from an external configuration file, the naming convention 
should be calculated by a helper function, implemented in a shared helper class.

### Schema file of the wrapper

If the wrapper class is for a specific template ( in this case for the CFT Folders template), a Schema file can be
used for the following:

 - Importing the target template makes the YAML easier and explicitly states the template dependency
 - Copying the required and optional property definition from the target template enforces the property validation in an earlier 
 stage. ( Unfortunately referencing to another Schema file is not possible today.)
 - Comments in the Schema file explains the usage and the purpose of it.

 ## Generic wrapper

 Using a generic wrapper fits into the concept of hierarchical configuration management when the configuration properties
 of the deployment are coming from multiple external files, not only the starting YAML. (See ../../hierarchical_configuration) 
 The generic wrapper is able to inject the context aware properties and pass them to the target template which is defined in
 the starting YAML.

 A nice trick is to import the target template in the YAML file and name it as "target-template.py", this makes you able to
 use the same wrapper template with any YAML/Target template combination.