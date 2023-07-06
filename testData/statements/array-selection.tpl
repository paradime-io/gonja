{% if foobar[:6] == 'barbaz' %}
    foo
{% else %}
    bar
{% endif %}

{% if foobar[0:] == 'barbaz' %}
    foo
{% else %}
    bar
{% endif %}

{% if foobar[:] == 'barbaz' %}
    foo
{% else %}
    bar
{% endif %}

{% if foobar[0:6] == 'barbaz' %}
    foo
{% else %}
    bar
{% endif %}