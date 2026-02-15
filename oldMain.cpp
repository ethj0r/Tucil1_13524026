#include <iostream>
#include <vector>
#include <cmath>
#include <chrono>
using namespace std;

vector<vector<int>> colors;


bool isSafe(vector<vector<int>> &mat, int row, int col) {
    int n = mat.size();

    //check row col
    for (int x=0; x<n; x++) {
        if (mat[row][x] ==1) return false;
    }
    for (int x=0; x<n; x++) {
        if (mat[x][col] == 1) return false;
    }

    //check diag
    for (int i=0; i<n; i++) {
        for (int j=0; j<n; j++) {
            if (mat[i][j] == 1 && abs(i-row) == abs(j-col)) return false;
        }
    }

    //check color
    int currColor = colors[row][col];
    for (int i=0; i<n; i++) {
        for (int j=0; j<n; j++) {
            if (mat[i][j]==1 && colors[i][j]==currColor) return false;
        }
    }
    
    return true;
}

bool solveQueensRec(vector<vector<int>> &mat, int row, int col) {
    int n = mat.size();

    if (row==n-1 && col==n) return true;
    if (col==n) {
        row++;
        col=0;
    }
    if (mat[row][col] != 0) return solveQueensRec(mat, row, col+1);
    if (isSafe(mat, row, col)) {
        mat[row][col] = 1;
        if (solveQueensRec(mat, row, col+1)) return true;
        mat[row][col] = 0;
    }

    return solveQueensRec(mat, row, col+1);
}

void solveQueens(vector<vector<int>> &mat) {
    solveQueensRec(mat, 0, 0);
}

int main() {
    int n = 8;
    vector<vector<int>> mat(n, vector<int>(n,0));
    colors = {
        {0, 1, 2, 3, 0, 1, 2, 3},
        {1, 2, 3, 0, 1, 2, 3, 0},
        {2, 3, 0, 1, 2, 3, 0, 1},
        {3, 0, 1, 2, 3, 0, 1, 2},
        {0, 1, 2, 3, 0, 1, 2, 3},
        {1, 2, 3, 0, 1, 2, 3, 0},
        {2, 3, 0, 1, 2, 3, 0, 1},
        {3, 0, 1, 2, 3, 0, 1, 2}
    };
    auto start = chrono::high_resolution_clock::now();
    solveQueens(mat);
    auto end = chrono::high_resolution_clock::now();
    chrono::duration<double> lama=end-start;

    for (int i=0; i<n; i++) {
        for (int j=0; j<n; j++) {
            cout << (mat[i][j]==1 ? "Q" : ".") << " ";
        }
        cout << "\n";
    }
    cout << "Execution time: " << lama.count() << " s\n";
    
    return 0;
}