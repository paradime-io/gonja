{% set a={"b":{"key":"value"}} %}
{% set b="b" %}
{% set c={"d":"key"} %}
{% set d="d" %}

{{ a.b[c.d] }}

{% set l=[1,2,3] %}
this here {{ l[l[l[0]]] }} should be 3