#version 330

in vec2 fragTexCoord;

uniform sampler2D texture0;

uniform float rewind = 0.0;

// black and white
vec4 blackWhite(vec4 color) {
    color.rgb = vec3((color.r + color.g + color.b)/3.0);
    return color;
}

void main()
{
    vec4 color = texture2D(texture0, fragTexCoord);
    if(rewind == 1.0) {
        color = blackWhite(color);
    }

    gl_FragColor = color;
}