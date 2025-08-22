number age = 18
text name = "Alice"

if age >= 18 then
    print name + " is an adult"
else
    print name + " is a minor"
end

number score = 85
if score >= 90 then
    print "Grade: A"
else
    if score >= 80 then
        print "Grade: B"
    else
        if score >= 70 then
            print "Grade: C"
        else
            print "Grade: F"
        end
    end
end
