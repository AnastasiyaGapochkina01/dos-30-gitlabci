import app

def test_add():
    assert app.add(5, 6) == 11
    assert app.add(-5, 10) == 5
