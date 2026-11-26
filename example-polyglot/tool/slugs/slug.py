"""A title's slug, as the service makes it."""

import os
import re

from slugs import env


# charter:realizes SW-261126-4NQQ
def slug(title, separator="-"):
    """The title's letters and digits in lower case, each run joined by one
    separator."""
    words = re.findall(r"[^\W_]+", title)
    return separator.join(words).lower()


# charter:realizes SW-261126-X6HB
def separator():
    """The separator the tool joins words with."""
    return os.environ.get(env.SEPARATOR) or "-"
