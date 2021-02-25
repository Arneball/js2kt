    go install github.com/Arneball/js2kt
# 
    cat <<EOF | js2kt
    {"negra": [{"s": 3, "testdata": "abc", "time": "2021-02-25T13:45:24+00:00", "bool": true, "uuid": "f93e2204-91c7-4a86-bc6f-f65e0b892563"}]}
    EOF

# 

Output:

    data class ChangeMe0(
        val negra: List<ChangeMe1>,
    )

    data class ChangeMe1(
        val s: Long,
        val testdata: String,
        val time: java.util.Date,
        val bool: Boolean,
        val uuid: java.util.UUID,
    )
