#version 330

in vec2 fragTexCoord;

uniform sampler2D texture0;

uniform vec2 objectPosBottomLeft;
uniform vec2 objectSize;

uniform float lightPosSize = 1;
uniform vec2 lightPos[10];

const float lightMaxDistance = 2000.0;
const float minBrightness = 0.01;
const float maxBrightness = 1.0;

void main()
{
    // Сэмплирование цвета пикселя из текстуры
    vec4 color = texture2D(texture0, fragTexCoord);

    vec2 posInPixel = fragTexCoord * objectSize;
    vec2 worldPos = objectPosBottomLeft + posInPixel;

    float brightness = 0.0;

    for (int i = 0; i < int(lightPosSize); i++) {
        vec2 lPos = lightPos[i];

        float distance = distance(worldPos, lPos);
        float diff = distance/lightMaxDistance;
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

