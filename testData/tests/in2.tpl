{% macro is_not_in(ls, value, if_true, if_false) -%}
    {%- if value|lower not in ls -%}
        {{- if_true -}}
    {%- else %}
        {{- if_false -}}
    {%- endif %}
{%- endmacro -%}


# contains

result: {{ is_not_in(["a","b"], "A", "does not contain", "contains") }}

# does not contain

result: {{ is_not_in(["a","b"], "C", "does not contain", "contains") }}