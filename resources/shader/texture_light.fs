#version 330

in vec2 fragTexCoord;

uniform sampler2D texture0;

uniform vec2 objectPosCenter;

uniform float lightPosSize = 1;
uniform vec2 lightPos[10];
uniform float lightMaxDistance[10];

const float minBrightness = 0.01;
const float maxBrightness = 1.0;

void main()
{
    // Сэмплирование цвета пикселя из текстуры
    vec4 color = texture2D(texture0, fragTexCoord);

    float brightness = 0.0;

    for (int i = 0; i < int(lightPosSize); i++) {
        vec2 lPos = lightPos[i];

        float distance = distance(objectPosCenter, lPos);
        float diff = distance/lightMaxDistance[i];
        if(diff >= 1) {
            brightness += minBrightness;
        } else {
            brightness += maxBrightness - diff;
        }
    }

    if (brightness > maxBrightness) {
        brightness = maxBrightness;
    }

    if (brightness < minBrightness) {
        brightness = minBrightness;
    }

    // Изменение яркости цвета пикселя
    color.rgb *= brightness;

    gl_FragColor = color;
}

