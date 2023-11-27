"""
Everything rots into incomprehensible complexity over time. One example of this is how pytest
detects the code to be tested.

> pytest will find foo/bar/tests/test_foo.py and realize it is part of a package given that
> there’s an __init__.py file in the same folder. It will then search upwards until it can find
> the last folder which still contains an __init__.py file in order to find the package root
> (in this case foo/).

https://docs.pytest.org/en/7.1.x/explanation/pythonpath.html#test-modules-conftest-py-files-inside-packages
"""