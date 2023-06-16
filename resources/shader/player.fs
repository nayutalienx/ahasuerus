#version 330

in vec2 fragTexCoord;

uniform sampler2D texture0;

uniform vec2 objectPosCenter;

uniform float lightPosSize = 1;
uniform vec2 lightPos[10];
uniform float lightMaxDistance[10];

uniform float playerWidth = 200;
uniform float playerHeight = 400;

const float minBrightness = 0.01;
const float maxBrightness = 1.0;

// blur

float offset[3] = float[](0.0, 1.3, 3.2);
float weight[3] = float[](0.2, 0.3, 0.07);

const float blurScale = 0.8;

vec4 blur(vec4 color) {
    
    // Texel color fetching from texture sampler
    vec3 texelColor = color.rgb*weight[0];

    float width = playerWidth*blurScale;

    vec2 texCoord = fragTexCoord;

    for (int i = 1; i < 3; i++)
    {
        texelColor += texture2D(texture0, texCoord + vec2(offset[i])/width, 0.0).rgb*weight[i];
        texelColor += texture2D(texture0, texCoord - vec2(offset[i])/width, 0.0).rgb*weight[i];
    }

    return vec4(texelColor, color.w);
}

// tail flow

// vec4 tailFlow(vec4 color) {
//     vec4 result = color;
//     if (color.a < 0.1) {
        
//         vec2 offset = vec2(0.05, 0.0);
//         vec2 rightPoint = fragTexCoord + offset;
//         if(rightPoint.x > 1.0) {
//             rightPoint.x = 1.0;
//         }
//         //vec2 leftPoint = fragTexCoord - offset;

//         vec4 rightTexel = texture2D(texture0, rightPoint);
//         //vec4 leftTexel = texture2D(texture0, leftPoint);

//         result = vec4(rightPoint.x, 0.0, 0.0, color.a+rightTexel.a);
//     }
//     return result;
// }

// main
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

    if (brightness < 0.1) {
        brightness = 0;
    }

    // Изменение яркости цвета пикселя
    color.rgb *= brightness;

    gl_FragColor = color;
}
