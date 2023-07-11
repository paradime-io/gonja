{% set default = "default" if not default %}
default: '{{ default }}'

{{ '\n' }}

{%- set my_override = my_override if my_override is not none else "fallback" -%}
first: '{{ my_override }}'

{{ '\n' }}

{%- set my_override = "already_set" -%}
{%- set my_override = my_override if my_override is not none else "fallback" -%}
second: '{{ my_override }}'