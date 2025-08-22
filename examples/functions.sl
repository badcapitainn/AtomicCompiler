function greet(text name)
    print "Hello, " + name + "!"
end

function add(number a, number b)
    number result = a + b
    print "The sum of " + a + " and " + b + " is " + result
end

function factorial(number n)
    if n <= 1 then
        print "Factorial of " + n + " is 1"
    else
        number result = 1
        loop i from 1 to n
            result = result * i
        end
        print "Factorial of " + n + " is " + result
    end
end

greet("World")
add(5, 3)
factorial(5)
