#version 330


in vec2 fragTexCoord;
uniform sampler2D texture0;
const float outline = 0.2;

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

    gl_FragColor = sampled;
}