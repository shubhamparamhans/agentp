````markdown
# Type Casting Fix for UUID and Special Types

**Date**: January 26, 2026  
**Issue**: Group by queries failing with UUID types from auto-generated models.json
**Root Cause**: Hardcoded type support only for integer/decimal; PostgreSQL type inference failing for UUID  
**Status**: ✅ **FIXED**

---

... (full content copied from TYPE_CASTING_FIX.md)

````
