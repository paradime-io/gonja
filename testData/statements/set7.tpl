{%- set invalid_escape_sequence = "\/" %}

invalid_escape_sequence; '{{ invalid_escape_sequence }}'

{%- set valid_escape_sequence = "\n" %}

valid_escape_sequence; '{{ valid_escape_sequence }}'