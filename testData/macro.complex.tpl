{%- macro lengther() -%}
    the tuple {{ varargs }} has length: {{ varargs | length }}
{%- endmacro -%}

{%- macro how_many_arguments() -%}
    {{ lengther(*varargs) }}
{%- endmacro -%}

{{ how_many_arguments(1) }}
{{ how_many_arguments(1,2) }}
{{ how_many_arguments(1,2,3) }}
{{ how_many_arguments(1,2,3,"a") }}

{{ '\n' }}

{%- macro itemiser() -%}
{% for key, value in kwargs.items() %} {{key}} -> {{value}} {% endfor %} ({{ kwargs | length }})
{%- endmacro -%}

{%- macro describe_my_kwargs() -%}
the dict {{ kwargs }} has length: {{ kwargs | length }}
it has these elements: {{ itemiser(**kwargs) }}
{%- endmacro -%}

{{ describe_my_kwargs(a=1) }}
{{ describe_my_kwargs(a=1,b=2,c="d") }}
