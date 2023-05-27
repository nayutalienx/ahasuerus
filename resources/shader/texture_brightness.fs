#version 330

in vec2 fragTexCoord;

uniform sampler2D texture0;

uniform vec2 objectPosBottomLeft;
uniform vec2 objectSize;

uniform vec2 lightPos;

const float lightMaxDistance = 2000.0;
const float minBrightness = 0.01;
const float maxBrightness = 1.0;

void main()
{
    // Сэмплирование цвета пикселя из текстуры
    vec4 color = texture2D(texture0, fragTexCoord);

    vec2 posInPixel = fragTexCoord * objectSize;
    vec2 worldPos = objectPosBottomLeft + posInPixel;

    float distance = distance(worldPos, lightPos);

    float diff = distance/lightMaxDistance;

    float brightness = 0.0;

    if(diff >= 1) {
        brightness = minBrightness;
    } else {
        brightness = maxBrightness - diff;
    }

    // Изменение яркости цвета пикселя
    color.rgb *= brightness;

    gl_FragColor = color;
}

