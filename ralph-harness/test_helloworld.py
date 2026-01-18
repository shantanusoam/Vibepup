import unittest
import helloworld
from io import StringIO
import sys


class TestHelloWorld(unittest.TestCase):
    def test_greeting_function(self):
        """Test the string return value directly."""
        self.assertEqual(helloworld.get_greeting(), "Hello from Ralph")


if __name__ == "__main__":
    unittest.main()
