{%- set nested = [
   {
      "key1":"value1",
      "key2":[
         {
            "key2.1":"value.2.1",
            "key.2.2":"value.2.2"
         }
      ]
   }
]
-%}

{{ nested }}