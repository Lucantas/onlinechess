function moveRange(piece, position){
    switch(piece){
        case "pawn":
            movePawn(position)
            break;
        case "rook":
            moveRook(position)
            break;
        case "knight":
            moveKnight(position)
            break;
        case "bishop":
            moveBishop(position)
            break;
        case "queen":
            moveQueen(position)
            break;
        case "king":
            moveking(position)
            break;
    }
}