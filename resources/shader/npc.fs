#version 330


in vec2 fragTexCoord;

uniform sampler2D texture0;
uniform float rewind = 0.0;

const float outline = 0.2;

// black and white
vec4 blackWhite(vec4 color) {
    color.rgb = vec3((color.r + color.g + color.b)/3.0);
    return color;
}

void main()
{
    vec4 sampled = texture2D(texture0, fragTexCoord);

    if (fragTexCoord.x < outline) {
        sampled.a = fragTexCoord.x/outline;
    }

    if (fragTexCoord.y > 1-outline) {
        sampled.a = (1.0 - fragTexCoord.y)/outline;

        if (fragTexCoord.x < outline) {
            sampled.a = min(fragTexCoord.x, 1-fragTexCoord.y)/outline;
        }
    }

    if (rewind == 1.0) {
        sampled = blackWhite(sampled);
    }

    gl_FragColor = sampled;
}