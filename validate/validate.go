package validate

/* ValidateName returns whether a module/variable/function/type name is valid.
 * The name must start with an alphabetical character, and contain only alpha-
 * numeric characters. Symbols such as dashes and underscores are not
 * considered valid.
 */
func ValidateName (name string) (valid bool) {
        runes := []rune(name)
        // symbols must be at least two letters, and start with an alphabetic
        // character.
        if len(runes) < 2 { return false }
        lower := runes[0] < 'a' || runes[0] > 'z'
        upper := runes[0] < 'Z' || runes[0] > 'Z'
        if !lower && !upper { return false }
        
        for _, ch := range runes {
                // keep going if the character is valid
                if ch >= 'a' && ch <= 'z' { continue }
                if ch >= 'A' && ch <= 'Z' { continue }
                if ch >= '0' && ch <= '9' { continue }
                
                // otherwise, return false
                return false
        }

        return true
}

/* ValidatePermission returns whether a permission is valid or not. For it to be
 * valid, it must be two runes in length and only contain the characters rwn.
 */
func ValidatePermission (permission string) (valid bool) {
        runes := []rune(permission)

        if len(runes) != 2 { return false }

        self  := runes[0]
        other := runes[1]

        if self  != 'r' && self  != 'w' && self  != 'n' { return false }
        if other != 'r' && other != 'w' && other != 'n' { return false }

        return true
}
