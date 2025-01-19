/*

## Men's Sizing Guide (Ages 18-60)
| Size | Height (cm) | Weight (kg) | Abdomen         | Chest         | Fit Preference   |
|------|-------------|-------------|-----------------|---------------|------------------|
| XS   | 155-165     | 45-55       | Flat            | Narrow        | Any              |
| S    | 160-170     | 55-65       | Flat/Medium     | Narrow/Medium | Any              |
| M    | 165-175     | 65-75       | Any             | Any           | Any              |
| L    | 170-180     | 75-85       | Medium/Building | Medium/Wide   | Normal/Loose     |
| XL   | 175-185     | 85-95       | Building        | Wide          | Normal/Loose     |
| XXL  | 180-190     | 95-105      | Building        | Wide          | Loose/Very Loose |

## Women's Sizing Guide (Ages 18-60)
| Size | Height (cm) | Weight (kg) | Abdomen         | Chest         | Fit Preference   |
|------|-------------|-------------|-----------------|---------------|------------------|
| XS   | 150-160     | 40-50       | Flat            | Narrow        | Any              |
| S    | 155-165     | 45-55       | Flat/Medium     | Narrow/Medium | Any              |
| M    | 160-170     | 50-65       | Any             | Any           | Any              |
| L    | 165-175     | 65-75       | Medium/Building | Medium/Wide   | Normal/Loose     |
| XL   | 170-180     | 75-85       | Building        | Wide          | Normal/Loose     |
| XXL  | 175-185     | 85-95       | Building        | Wide          | Loose/Very Loose |

## Fit Adjustment Guidelines:

1. Age Considerations:
   - Ages 18-30: Can stay with chart recommendation
   - Ages 30-45: Consider going up one size if between sizes
   - Ages 45+: May prefer going up one size for comfort

2. Body Shape Adjustments:
   - Flat Abdomen: Can choose smaller size if between sizes
   - Medium Abdomen: Stay with chart recommendation
   - Building Abdomen: Consider going up one size

3. Fit Preference Adjustments:
   - Very Fitted: Go down one size
   - Fitted: Use chart size
   - Normal: Use chart size
   - Loose: Go up one size
   - Very Loose: Go up two sizes

4. Special Considerations:
   - Athletic Build: May need to size up for shoulders/chest
   - Height-Weight Proportions: Prioritize weight over height if conflicting
   - Between Sizes: Consider fit preference and body shape

*/

function calculateSize(data) {
    let baseSize = '';
    const { gender, height, weight } = data;

    if (gender === 'male') {
      if (height < 165 && weight < 55) baseSize = 'XS';
      else if (height < 170 && weight < 65) baseSize = 'S';
      else if (height < 175 && weight < 75) baseSize = 'M';
      else if (height < 180 && weight < 85) baseSize = 'L';
      else if (height < 185 && weight < 95) baseSize = 'XL';
      else baseSize = 'XXL';
    } else {
      if (height < 160 && weight < 50) baseSize = 'XS';
      else if (height < 165 && weight < 55) baseSize = 'S';
      else if (height < 170 && weight < 65) baseSize = 'M';
      else if (height < 175 && weight < 75) baseSize = 'L';
      else if (height < 180 && weight < 85) baseSize = 'XL';
      else baseSize = 'XXL';
    }

    let sizeAdjustment = 0;
    if (data.age > 30 && data.age <= 45) sizeAdjustment += 0.5;
    if (data.age > 45) sizeAdjustment += 1;
    if (data.abdomen === 'building') sizeAdjustment += 1;
    if (data.fitPreference === 'very_fitted') sizeAdjustment -= 1;
    if (data.fitPreference === 'loose') sizeAdjustment += 1;
    if (data.fitPreference === 'very_loose') sizeAdjustment += 2;

    const sizeMap = { 'XS': 0, 'S': 1, 'M': 2, 'L': 3, 'XL': 4, 'XXL': 5 };
    const sizeArray = ['XS', 'S', 'M', 'L', 'XL', 'XXL'];
    const finalSizeIndex = Math.min(Math.max(sizeMap[baseSize] + Math.round(sizeAdjustment), 0), 5);

    return { recommendedSize: sizeArray[finalSizeIndex], baseSize };
}

// const formData = {
//     gender: "male",
//     age: 19,
//     height: 170,
//     weight: 50,
//     abdomen: "medium",
//     chest: "medium",
//     fitPreference: "normal",
// };

// const result = calculateSize(formData);
// console.log(result)