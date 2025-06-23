"""
SortingHat models for resume scoring service.

This module contains Pydantic models that correspond to the TypeSpec definitions
in sortinghat.tsp, allowing for easy import and use in Python code.

This module automatically exports all public classes from all Python files
in this directory, so you don't need to manually update imports when adding new types.
"""

# Import all exports from sortinghat module
# This automatically picks up everything in sortinghat.py's __all__ list
from .sortinghat import * 