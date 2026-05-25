## 1. THI Helpers in index.html

- [x] 1.1 Add `thi(tempC, rh)` function computing Thom's Discomfort Index
- [x] 1.2 Add `thiLevel(thi)` function returning `{color, label}` for the six comfort bands

## 2. Update render() in index.html

- [x] 2.1 Compute THI from each sensor's latest reading inside `render()`
- [x] 2.2 Apply THI border color to both metric cards (temperature + humidity)
- [x] 2.3 Append `· <label>` in THI color to the timestamp row
- [x] 2.4 Add `title` attribute to the label `<span>` showing the numeric THI value (1 decimal place)
