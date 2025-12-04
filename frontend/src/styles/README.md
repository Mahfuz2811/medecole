# Design System Usage Guide

## 🎨 Design Tokens Implementation

This design system solves the following issues:

### ✅ **Fixed Issues:**

1. **Inconsistent Border Radius** - Now using systematic `designTokens.radius.*`
2. **Hard-coded Values** - Replaced with reusable tokens
3. **Redundant Styling** - Eliminated duplicate classes
4. **No System** - Created consistent design language

### 🚀 **How to Use:**

```tsx
import { designTokens, layouts, cn } from "@/styles/design-tokens";

// ✅ Use design tokens instead of hard-coded classes
<div className={designTokens.components.card.interactive}>

// ✅ Combine multiple tokens
<div className={cn(layouts.flexBetween, designTokens.spacing.md)}>

// ✅ Use avatar variants
<div className={designTokens.components.avatar.md}>
```

### 📦 **Component Patterns:**

```tsx
// Cards
designTokens.components.card.base; // Basic card
designTokens.components.card.interactive; // Hover effects
designTokens.components.card.bordered; // With border

// Layout
layouts.container; // Responsive container
layouts.flexBetween; // Space between layout
layouts.flexCenter; // Centered layout
layouts.stickyHeader; // Fixed header positioning

// Colors
designTokens.colors.background.primary; // bg-white
designTokens.colors.text.accent; // text-blue-600
designTokens.colors.background.accent; // bg-blue-50
```

### 🔧 **Migration Strategy:**

**Before:**

```tsx
<div className="bg-white shadow-sm p-4 rounded-lg hover:shadow-md transition">
```

**After:**

```tsx
<div className={designTokens.components.card.interactive}>
```

### 💡 **Benefits:**

1. **Consistency** - All components use same design language
2. **Maintainability** - Change tokens in one place
3. **Type Safety** - TypeScript autocompletion
4. **Performance** - Optimized class combinations
5. **Scalability** - Easy to extend and modify

### 🎯 **Next Steps:**

1. Apply to remaining components (ScheduleCard, RankingCard, etc.)
2. Add more component variants as needed
3. Extend color system for themes
4. Add animation/transition tokens
