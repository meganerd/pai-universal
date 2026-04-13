# WISDOM.md

## Engineering Wisdom from Experience

### Principles distilled from building high-performance software

---

### 1. Software Has Weight
> "You never fully recover from [programming on a 6502]. You stop seeing software as this soft, fluffy thing made out of abstractions, and start seeing it as having weight, mass, drag, and friction."

Every line has a cost. Every allocation leaves footprints. Every dependency is a roommate that eats your food and never pays rent.

---

### 2. Small Was Fast, Fast Mattered
> "Small was fast, and fast mattered."

Small is fast. Fast is responsive. Responsive builds trust. Size and performance are features, not afterthoughts.

---

### 3. Tools Invoked in Crisis Must Be Crisp
> "Task Manager is the sort of program that you invoke when things are already going wrong... It has to be there now, and it has to feel crisp."

When the user needs your tool, it can't arrive fashionably late, staggering under the weight of its own dependencies. Show up like Victor, always wearing the tuxedo.

---

### 4. Never Lie About the Data
> "Never lie about the data."

Accuracy over ease. "Good enough" is not acceptable. If you're showing numbers, they'd better be right.

---

### 5. Eat Your Own Dog Food
> "Eat your own dog food"

Ship builds. Test in production-like conditions. Take quality personally until the fix ships.

---

### 6. Give Power, Not Comfort
> "If the user needs a chisel, don't give them a Nerf bat"

Don't dumb down tools for perceived ease of use. Trust the user with power tools.

---

### 7. Build for the Hardware You Have
> "Build for the hardware you have, not hope"

Don't assume faster machines will solve your performance problems. Optimize for reality.

---

### 8. Small, Accurate, Robust
These are the priorities, not "shipping features".

---

### 9. Three Priorities for Task Manager
1. Dynamic resizing with no flicker
2. Keep it small
3. Robustness

---

### 10. Every Allocation Leaves Footprints
Memory allocations aren't free. They have cost in time, space, and cache behavior.

---

### 11. Dependencies Are Roommates
Every dependency you add eats your resources and never pays rent. Choose carefully.

---

### 12. Complexity Compounds
Each "layer of comfort" and "future-proofing" adds weight. The next person maintaining your code will thank you for restraint.

---

### 13. Fast Is a Feature
Performance is a feature, not something to add "later". Later never comes.

---

*Source: Dave Plummer - Windows Task Manager creator, Dave's Garage*
*Extracted from: "Why the Original Task Manager Was Under 80K and Insanely Fast"